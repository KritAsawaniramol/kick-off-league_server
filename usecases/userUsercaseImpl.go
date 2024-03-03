package usecases

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"kickoff-league.com/config"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type userUsecaseImpl struct {
	userrepository repositories.Userrepository
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageProfile(normalUserID uint) error {
	normalUser := entities.NormalUsers{}
	normalUser.ID = normalUserID
	if err := u.userrepository.UpdateSelectedFields(normalUser, "ImageProfilePath", &entities.NormalUsers{ImageProfilePath: ""}); err != nil {
		return err
	}
	return nil
}

// GetCompatitions implements UserUsecase.
func (*userUsecaseImpl) GetCompatitions(in *model.GetCompatitionsReq) ([]model.Compatition, error) {
	return nil, nil
}

// CreateCompatition implements UserUsecase.
func (u *userUsecaseImpl) CreateCompatition(in *model.CreateCompatition) error {
	if err := u.userrepository.InsertCompatition(&entities.Compatitions{
		Name:         in.Name,
		Format:       entities.CompetitionFormat(in.Format),
		OrganizersID: in.OrganizerID,
		StartDate:    in.StartDate,
		EndDate:      in.EndDate,
		AgeOver:      in.AgeOver,
		AgeUnder:     in.AgeUnder,
		Sex:          entities.SexType(in.Sex),
		Description:  in.Description,
		Status:       "creating",
	}); err != nil {
		return err
	}
	return nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// GetMyPenddingAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) GetMyPenddingAddMemberRequest(userID uint) ([]model.AddMemberRequest, error) {
	addMemberRequestSearch := &entities.AddMemberRequests{
		ReceiverID: userID,
		Status:     "pending",
	}
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return nil, err
	}

	addMemberRequestModelList := []model.AddMemberRequest{}
	for _, v := range addMemberRequestList {
		addMemberRequestModelList = append(addMemberRequestModelList, model.AddMemberRequest{
			ID:       v.ID,
			TeamID:   v.TeamsID,
			TeamName: v.Teams.Name,
			Role:     v.Role,
			Status:   v.Status,
		})
	}

	return addMemberRequestModelList, nil
}

// GetTeamMember implements UserUsecase.
func (u *userUsecaseImpl) GetTeamMembers(id uint) (*model.Team, error) {
	teamsMembers, err := u.userrepository.GetTeamMembersByTeamID(id, "id", false, -1, -1)
	if err != nil {
		return nil, err
	}
	util.PrintObjInJson(teamsMembers)
	return nil, nil
}

// GetTeamWithMemberAndCompatitionByID implements UserUsecase.
func (u *userUsecaseImpl) GetTeamWithMemberAndCompatitionByID(id uint) (*model.Team, error) {
	t := &entities.Teams{}
	t.ID = id
	selectedTeams, err := u.userrepository.GetTeamWithMemberAndCompatitionByID(id)
	if err != nil {
		return &model.Team{}, err
	}
	util.PrintObjInJson(selectedTeams)
	memberList := []model.Member{}
	for _, member := range selectedTeams.TeamsMembers {
		memberList = append(memberList, model.Member{
			FirstNameThai: member.NormalUsers.FirstNameThai,
			LastNameThai:  member.NormalUsers.LastNameThai,
			FirstNameEng:  member.NormalUsers.FirstNameEng,
			LastNameEng:   member.NormalUsers.LastNameEng,
			Position:      member.NormalUsers.Position,
			Sex:           member.NormalUsers.Sex,
			Role:          member.Role,
		})
	}

	compatition_model := []model.CompatitionBasicInfo{}

	for _, v := range selectedTeams.Compatitions {
		compatition_model = append(compatition_model, model.CompatitionBasicInfo{
			ID:                v.ID,
			Name:              v.Name,
			Format:            model.CompetitionFormat(v.Format),
			OrganizerID:       v.OrganizersID,
			StartDate:         v.StartDate,
			EndDate:           v.EndDate,
			RegisterStartDate: v.RegisterStartDate,
			RegisterEndDate:   v.RegisterEndDate,
			ApplicationFee:    v.ApplicationFee,
			AgeOver:           v.AgeOver,
			AgeUnder:          v.AgeUnder,
			Sex:               model.SexType(v.Sex),
			FieldSurface:      model.FieldSurfaces(v.FieldSurface),
			Description:       v.Description,
			Status:            model.CompetitionStatus(v.Status),
			NumberOfTeam:      v.NumberOfTeam,
		})
	}

	return &model.Team{
		ID:           selectedTeams.ID,
		Name:         selectedTeams.Name,
		OwnerID:      selectedTeams.OwnerID,
		Members:      memberList,
		Compatitions: compatition_model,
		Description:  selectedTeams.Description,
	}, nil
}

// GetTeamList implements UserUsecase.
func (u *userUsecaseImpl) GetTeams(in *model.GetTeamsReq) ([]model.TeamList, error) {

	team := entities.Teams{
		// 0 is select all id
		OwnerID: in.NormalUserID,
	}
	team.ID = in.TeamID

	limit := int(in.PageSize)
	if limit <= 0 {
		limit = -1
	}

	offset := int(in.PageSize * in.Page)
	if offset <= 0 {
		offset = -1
	}

	if in.Ordering == "" {
		in.Ordering = "id"
	}
	util.PrintObjInJson(team)

	teams, err := u.userrepository.GetTeams(&team, strings.Trim(in.Ordering, " "), in.Decs, limit, offset)
	if err != nil {
		return []model.TeamList{}, err
	}

	teamList := []model.TeamList{}
	for _, team := range teams {
		teamList = append(teamList, model.TeamList{
			ID:             team.ID,
			Name:           team.Name,
			Description:    team.Description,
			NumberOfMember: uint(u.userrepository.GetNumberOfTeamsMember(team.ID)),
		})
	}

	return teamList, nil
}

// GetMyTeam implements UserUsecase.
func (*userUsecaseImpl) GetMyTeam() ([]model.Team, error) {
	panic("unimplemented")
}

// IgnoreAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) IgnoreAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]
	if addMemberRequest.ReceiverID != userID {
		return errors.New("UserID does not match ReceiverID")
	}
	if err := u.userrepository.UpdateAddMemberRequestStatusAndSoftDelete(addMemberRequest, "ignore"); err != nil {
		return err
	}
	return nil
}

// AcceptAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) AcceptAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequestSearch := &entities.AddMemberRequests{}
	addMemberRequestSearch.ID = inReqID
	addMemberRequestList, err := u.userrepository.GetAddMemberRequestByID(addMemberRequestSearch)
	if err != nil {
		return err
	}
	addMemberRequest := &addMemberRequestList[0]

	if addMemberRequest.ReceiverID != userID {
		log.Printf("ReceiverID: %d, userID: %d\n", addMemberRequest.ReceiverID, userID)
		return errors.New("UserID does not match ReceiverID")
	}

	if err := u.userrepository.UpdateAddMemberRequestStatusByID(inReqID, "accepted"); err != nil {
		return err
	}

	if err := u.userrepository.InsertTeamsMembers(&entities.TeamsMembers{
		TeamsID:       addMemberRequest.TeamsID,
		NormalUsersID: userID,
		Role:          addMemberRequest.Role,
	}); err != nil {
		return err
	}

	return nil
}

// CreateAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) SendAddMemberRequest(inAddMemberRequest *model.AddMemberRequest, inUserID uint) error {
	// Get Receiver normaluser
	receiver, err := u.userrepository.GetNormalUser(&entities.NormalUsers{
		Username: inAddMemberRequest.ReceiverUsername,
	})
	if err != nil {
		return err
	}

	team, err := u.userrepository.GetTeamWithMemberAndRequestSendByID(inAddMemberRequest.TeamID)
	if err != nil {
		return err
	}

	if team.OwnerID != inUserID {
		return errors.New("this user isn't owner's team")
	}

	for _, member := range team.TeamsMembers {
		if member.NormalUsers.UsersID == receiver.UsersID {
			return errors.New("this request receiver is already in team")
		}
	}

	for _, requestSend := range team.RequestSends {
		if requestSend.ReceiverID == receiver.UsersID && requestSend.Status == "pending" {
			return errors.New("team have already sent a request to this receiver")
		}
	}

	addMemberRequest := &entities.AddMemberRequests{
		TeamsID:    inAddMemberRequest.TeamID,
		ReceiverID: receiver.UsersID,
		Status:     "pending",
		Role:       inAddMemberRequest.Role,
	}

	// Create Request
	if err := u.userrepository.InsertAddMemberRequest(addMemberRequest); err != nil {
		return err
	}

	return nil

}

// UpdateNormalUser implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error {

	normalUser := &entities.NormalUsers{
		FirstNameThai:    inUpdateModel.FirstNameThai,
		LastNameThai:     inUpdateModel.LastNameThai,
		FirstNameEng:     inUpdateModel.FirstNameEng,
		LastNameEng:      inUpdateModel.LastNameEng,
		Born:             inUpdateModel.Born,
		Height:           inUpdateModel.Height,
		Weight:           inUpdateModel.Weight,
		Sex:              inUpdateModel.Sex,
		Position:         inUpdateModel.Position,
		Nationality:      inUpdateModel.Nationality,
		Description:      inUpdateModel.Description,
		Phone:            inUpdateModel.Phone,
		ImageProfilePath: inUpdateModel.ImageProfilePath,
		ImageCoverPath:   inUpdateModel.ImageCoverPath,
	}

	normalUser.ID = inNormalUserID

	util.PrintObjInJson(normalUser)

	if err := u.userrepository.UpdateNormalUser(normalUser); err != nil {
		return err
	}
	return nil
}

// CreateTeam implements UserUsecase.
func (u *userUsecaseImpl) CreateTeam(in *model.CreaetTeam) error {

	normalUser, err := u.userrepository.GetNormalUser(&entities.NormalUsers{
		UsersID: in.OwnerID,
	})
	if err != nil {
		return err
	}

	//check required data
	if normalUser.FirstNameThai == "" ||
		normalUser.LastNameThai == "" ||
		normalUser.FirstNameEng == "" ||
		normalUser.Born.IsZero() ||
		normalUser.Sex == "" ||
		normalUser.Nationality == "" ||
		normalUser.Phone == "" {
		return errors.New("no required normaluser data")
	}

	team := entities.Teams{
		Name:         in.Name,
		OwnerID:      in.OwnerID,
		Description:  in.Description,
		Compatitions: []entities.Compatitions{},
	}

	if err := u.userrepository.InsertTeam(&team); err != nil {
		return err
	}

	return nil
}

// UpdateNormalUserPhone implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUserPhone(inUserID uint, newPhone string) error {
	err := u.userrepository.UpdateNormalUserPhone(inUserID, newPhone)
	if err != nil {
		return err
	}
	return nil
}

func NewUserUsercaseImpl(
	userrepository repositories.Userrepository,
) UserUsecase {
	return &userUsecaseImpl{
		userrepository: userrepository,
	}
}

func (u *userUsecaseImpl) Login(in *model.LoginUser) (string, model.User, error) {

	if !isEmail(in.Email) {
		return "", model.User{}, errors.New("email is invalid")
	}

	//get user from email
	user, err := u.userrepository.GetUserByEmail(in.Email)
	if err != nil {
		return "", model.User{}, err
	}
	//compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password)); err != nil {
		return "", model.User{}, err
	}

	userModel := model.User{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	// pass = return jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	if user.Role == "normal" {
		normalUser, err := u.userrepository.GetNormalUserWithAddressByUserID(user.ID)
		if err != nil {
			return "", model.User{}, err
		}
		if normalUser.ImageProfilePath != "" {
			normalUser.ImageProfilePath = normalUser.ImageProfilePath[1:]
		}
		if normalUser.ImageCoverPath != "" {
			normalUser.ImageCoverPath = normalUser.ImageCoverPath[1:]
		}

		userModel.Datail = map[string]interface{}{

			"normal_user_info": model.NormalUserInfo{
				ID:               normalUser.ID,
				FirstNameThai:    normalUser.FirstNameThai,
				LastNameThai:     normalUser.LastNameThai,
				FirstNameEng:     normalUser.FirstNameEng,
				LastNameEng:      normalUser.LastNameEng,
				Born:             normalUser.Born,
				Phone:            normalUser.Phone,
				Height:           normalUser.Height,
				Weight:           normalUser.Weight,
				Sex:              normalUser.Sex,
				Position:         normalUser.Position,
				Nationality:      normalUser.Nationality,
				Description:      normalUser.Description,
				ImageProfilePath: normalUser.ImageProfilePath,
				ImageCoverPath:   normalUser.ImageCoverPath,
				Address: model.Address{
					HouseNumber: normalUser.Addresses.HouseNumber,
					Village:     normalUser.Addresses.Village,
					Subdistrict: normalUser.Addresses.Subdistrict,
					District:    normalUser.Addresses.District,
					PostalCode:  normalUser.Addresses.PostalCode,
					Country:     normalUser.Addresses.Country,
				},
			},
		}
		claims["normal_user_id"] = normalUser.ID
	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return "", model.User{}, err
		}
		userModel.Datail = map[string]interface{}{
			"organizer": model.OrganizersInfo{
				ID:          organizer.ID,
				Name:        organizer.Name,
				Phone:       organizer.Phone,
				Description: organizer.Description,
				Address: model.Address{
					HouseNumber: organizer.Addresses.HouseNumber,
					Village:     organizer.Addresses.Village,
					Subdistrict: organizer.Addresses.Subdistrict,
					District:    organizer.Addresses.District,
					PostalCode:  organizer.Addresses.PostalCode,
					Country:     organizer.Addresses.Country,
				},
			},
		}
		claims["organizer_id"] = organizer.ID
	}

	//bypass
	jwtSecretKey := config.GetConfig().JwtSecretKey

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", model.User{}, err
	}
	return t, userModel, nil
}
func (u *userUsecaseImpl) RegisterNormaluser(in *model.RegisterNormaluser) error {

	if !isEmail(in.Email) {
		return errors.New("email is invalid")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	in.Password = string(hashedPassword)

	user := &entities.Users{
		Email:    in.Email,
		Role:     "normal",
		Password: in.Password,
	}

	normalUser := &entities.NormalUsers{
		UsersID:  user.ID,
		Username: in.Username,
	}
	if err := u.userrepository.InsertUserWihtNormalUserAndAddress(normalUser, user); err != nil {
		return err
	}

	return nil
}

func isEmailAlreadyInUse(email string, u repositories.Userrepository) bool {
	if _, err := u.GetUserByEmail(email); err != nil {
		return false
	}
	return true
}

func isPhoneAlreadyInUse(phone string, u repositories.Userrepository) bool {
	if _, err := u.GetNormalUser(&entities.NormalUsers{
		Phone: phone,
	}); err != nil {
		return false
	}

	if _, err := u.GetOrganizer(&entities.Organizers{
		Phone: phone,
	}); err != nil {
		return false
	}

	return true
}

func (u *userUsecaseImpl) RegisterOrganizer(in *model.RegisterOrganizer) error {

	if !isEmail(in.Email) {
		return errors.New("email is invalid")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}

	if isPhoneAlreadyInUse(in.Phone, u.userrepository) {
		return errors.New("this phone is already in us")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	in.Password = string(hashedPassword)

	user := &entities.Users{
		Email:    in.Email,
		Role:     "organizer",
		Password: in.Password,
	}

	organizer := &entities.Organizers{
		UsersID: user.ID,
		Name:    in.OrganizerName,
		Phone:   in.Phone,
	}

	if err := u.userrepository.InsertUserWihtOrganizerAndAddress(organizer, user); err != nil {
		return err
	}

	return nil
}

// GetUsers implements UserUsecase.
func (u *userUsecaseImpl) GetUsers() ([]model.User, error) {
	users_entity, err := u.userrepository.GetUsers()
	if err != nil {
		return []model.User{}, err
	}
	users_model := []model.User{}
	for _, e := range users_entity {
		m := model.User{
			ID:        e.ID,
			Email:     e.Email,
			Role:      e.Role,
			CreatedAt: e.CreatedAt,
		}
		users_model = append(users_model, m)
	}

	return users_model, nil
}

// GetUser implements UserUsecase.
func (u *userUsecaseImpl) GetUser(in uint) (model.User, error) {
	user_entities, err := u.userrepository.GetUserByID(in)
	if err != nil {
		return model.User{}, err
	}

	user_model := model.User{
		ID:        user_entities.ID,
		Email:     user_entities.Email,
		Role:      user_entities.Role,
		CreatedAt: user_entities.CreatedAt,
	}

	return user_model, nil
}
