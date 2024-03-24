package usecases

import (
	"errors"
	"fmt"
	"math"
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

// GetNormalUserList implements UserUsecase.
func (u *userUsecaseImpl) GetNormalUserList() ([]model.NormalUserList, error) {
	normalUsers_entity, err := u.userrepository.GetNormalUsers(&entities.NormalUsers{})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	normalUserList := []model.NormalUserList{}
	for _, v := range normalUsers_entity {
		normalUserList = append(normalUserList, model.NormalUserList{
			ID:            v.ID,
			FirstNameThai: v.FirstNameThai,
			LastNameThai:  v.LastNameThai,
			FirstNameEng:  v.FirstNameEng,
			LastNameEng:   v.LastNameEng,
			Born:          v.Born,
			Height:        v.Height,
			Weight:        v.Weight,
			Sex:           v.Sex,
			Position:      v.Position,
			Nationality:   v.Nationality,
			Description:   v.Description,
		})
	}
	return normalUserList, nil
}

// UpdateUser implements UserUsecase.
func (*userUsecaseImpl) UpdateUser(in *model.User) error {
	panic("unimplemented")
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageProfile(userID uint, newImagePath string) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageProfilePath", &entities.NormalUsers{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) UpdateImageCover(userID uint, newImagePath string) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageProfilePath", &entities.NormalUsers{ImageProfilePath: newImagePath}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageProfile(userID uint) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageProfilePath", &entities.NormalUsers{ImageProfilePath: "./images/default/defaultProfile.jpg"}); err != nil {
		return err
	}
	return nil
}

// RemoveImageProfile implements UserUsecase.
func (u *userUsecaseImpl) RemoveImageCover(userID uint) error {
	user := &entities.Users{}
	user.ID = userID
	if err := u.userrepository.UpdateSelectedFields(user, "ImageCoverPath", &entities.NormalUsers{ImageCoverPath: "./images/default/defaultCover.jpg"}); err != nil {
		return err
	}
	return nil
}

// GetCompatitions implements UserUsecase.
func (*userUsecaseImpl) GetCompatitions(in *model.GetCompatitionsReq) ([]model.Compatition, error) {
	return nil, nil
}

// GetCompatition implements UserUsecase.
func (u *userUsecaseImpl) GetCompatition(in uint) (*model.GetCompatition, error) {
	compatitionEntity := &entities.Compatitions{}
	compatitionEntity.ID = in
	result, err := u.userrepository.GetCompatition(compatitionEntity)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	joinCode := []model.JoinCode{}

	for _, v := range result.JoinCode {
		joinCode = append(joinCode, model.JoinCode{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			Code:      v.Code,
			Status:    v.Status,
		})
	}

	temes := []model.Team{}

	for _, v := range result.Teams {
		members := []model.Member{}
		for _, member := range v.TeamsMembers {
			members = append(members, model.Member{
				ID:            member.ID,
				UsersID:       member.NormalUsers.UsersID,
				FirstNameThai: member.NormalUsers.FirstNameThai,
				LastNameThai:  member.NormalUsers.LastNameThai,
				FirstNameEng:  member.NormalUsers.FirstNameEng,
				LastNameEng:   member.NormalUsers.LastNameEng,
				Position:      member.NormalUsers.Position,
				Sex:           member.NormalUsers.Sex,
			})
		}

		temes = append(temes, model.Team{
			ID:          v.ID,
			Name:        v.Name,
			OwnerID:     v.OwnerID,
			Description: v.Description,
			Members:     members,
		})
	}

	matches := []model.Match{}

	for _, v := range result.Matches {
		goalRecords := []model.GoalRecord{}

		for _, goalRecord := range v.GoalRecords {
			goalRecords = append(goalRecords, model.GoalRecord{
				MatchesID:  goalRecord.MatchesID,
				TeamID:     goalRecord.TeamsID,
				PlayerID:   goalRecord.PlayerID,
				TimeScored: goalRecord.TimeScored,
			})
		}

		matches = append(matches, model.Match{
			ID:             v.ID,
			Index:          v.Index,
			DateTime:       v.DateTime,
			Team1ID:        v.Team1ID,
			Team2ID:        v.Team2ID,
			Team1Goals:     v.Team1Goals,
			Team2Goals:     v.Team2Goals,
			Round:          v.Round,
			NextMatchIndex: v.NextMatchIndex,
			NextMatchSlot:  v.NextMatchSlot,
			GoalRecords:    goalRecords,
			Result:         model.MatchesResult(v.Result),
		})
	}

	getCompatitionModel := model.GetCompatition{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		Name:      result.Name,
		Sport:     result.Sport,
		Format:    result.Format,
		Type:      model.CompetitionFormat(result.Type),
		OrganizerInfo: model.OrganizersInfo{
			ID:          result.OrganizersID,
			Name:        result.Organizers.Name,
			Phone:       result.Organizers.Phone,
			Description: result.Organizers.Description,
			Address: model.Address{
				HouseNumber: result.Organizers.Addresses.HouseNumber,
				Village:     result.Organizers.Addresses.Village,
				Subdistrict: result.Organizers.Addresses.Subdistrict,
				District:    result.Organizers.Addresses.District,
				PostalCode:  result.Organizers.Addresses.PostalCode,
				Country:     result.Organizers.Addresses.Country,
			},
			ImageProfilePath: result.Organizers.ImageProfilePath,
			ImageCoverPath:   result.Organizers.ImageCoverPath,
		},
		FieldSurface:    string(result.FieldSurface),
		ApplicationType: result.ApplicationType,
		Address: model.Address{
			HouseNumber: result.HouseNumber,
			Village:     result.Village,
			Subdistrict: result.Subdistrict,
			District:    result.District,
			PostalCode:  result.PostalCode,
			Country:     result.Country,
		},
		ImageBanner:          result.ImageBanner,
		StartDate:            result.StartDate,
		EndDate:              result.EndDate,
		JoinCode:             joinCode,
		Description:          result.Description,
		Rule:                 result.Rule,
		Prize:                result.Prize,
		ContractType:         result.ContractType,
		Contract:             result.Contract,
		AgeOver:              result.AgeOver,
		AgeUnder:             result.AgeUnder,
		Sex:                  model.SexType(result.Sex),
		Status:               string(result.Status),
		NumberOfTeam:         result.NumberOfTeam,
		NumOfPlayerInTeamMin: result.NumOfPlayerInTeamMin,
		NumOfPlayerInTeamMax: result.NumOfPlayerInTeamMax,
		Teams:                temes,
		NumOfRound:           result.NumOfRound,
		NumOfMatch:           result.NumOfMatch,
		Matches:              matches,
	}
	return &getCompatitionModel, nil
}

// CreateCompatition implements UserUsecase.
func (u *userUsecaseImpl) CreateCompatition(in *model.CreateCompatition) error {

	compatition := &entities.Compatitions{
		Name:                 in.Name,
		Sport:                in.Sport,
		Type:                 entities.CompetitionFormat(in.Type),
		Format:               in.Format,
		Description:          in.Description,
		Rule:                 in.Rule,
		Prize:                in.Prize,
		StartDate:            in.StartDate,
		EndDate:              in.EndDate,
		ApplicationType:      in.ApplicationType,
		ImageBanner:          in.ImageBanner,
		AgeOver:              in.AgeOver,
		AgeUnder:             in.AgeUnder,
		Sex:                  entities.SexType(in.Sex),
		NumberOfTeam:         in.NumberOfTeam,
		NumOfPlayerInTeamMin: in.NumOfPlayerInTeamMin,
		NumOfPlayerInTeamMax: in.NumOfPlayerInTeamMax,
		FieldSurface:         entities.FieldSurfaces(in.FieldSurface),
		OrganizersID:         in.OrganizerID,
		HouseNumber:          in.Address.HouseNumber,
		Village:              in.Address.Village,
		Subdistrict:          in.Address.Subdistrict,
		District:             in.Address.District,
		PostalCode:           in.Address.PostalCode,
		Country:              in.Address.Country,
		ContractType:         in.ContractType,
		Contract:             in.Contract,
		Status:               "Coming soon",
	}

	matchs := []entities.Matches{}
	numOfRound := 0
	if in.Type == "Tournament" {
		if checkNumberPowerOfTwo(int(in.NumberOfTeam)) != 0 {
			return errors.New("number of Team for create competition(tounament) is not power of 2")
		}
		if in.NumberOfTeam < 2 {
			return errors.New("number of Team have to morn than 1")
		}
		numOfRound = int(math.Log2(float64(in.NumberOfTeam)))
		fmt.Printf("in.NumberOfTeam: %v\n", in.NumberOfTeam)
		fmt.Printf("numOfRound: %v\n", numOfRound)
		count := 0
		for i := 0; i < numOfRound; i++ {
			round := numOfRound - i
			numOfMatchInRound := int(math.Pow(2, float64(round)) / 2)
			fmt.Printf("number of match in round %d: %d\n", round, numOfMatchInRound)
			for j := 0; j < int(numOfMatchInRound); j++ {
				match := entities.Matches{
					Round: fmt.Sprintf("Round %d", i+1),
				}
				if i != numOfRound-1 {
					if j%2 == 0 {
						match.NextMatchSlot = "Team1"
					} else {
						match.NextMatchSlot = "Team2"
					}
				}
				if i != 0 {
					fmt.Printf("i: %v\n", i)
					matchs[count].NextMatchIndex = len(matchs) + 1
					matchs[count+1].NextMatchIndex = len(matchs) + 1
					count += 2
				}
				match.Index = len(matchs) + 1
				matchs = append(matchs, match)
			}
		}
	} else if in.Type == "Round Robin" {
		numOfRound = int(in.NumberOfTeam - 1)
		numOfMatch := (int(in.NumberOfTeam) * numOfRound) / 2
		numOfMatchInRound := numOfMatch / numOfRound
		for i := 1; i <= int(numOfRound); i++ {
			for j := 0; j < int(numOfMatchInRound); j++ {
				matchs = append(matchs, entities.Matches{
					Round: fmt.Sprintf("Round %d", i),
					Index: len(matchs) + 1,
				})
			}
		}
	} else {
		return errors.New("undefined compatition type")
	}

	fmt.Printf("number of match %d\n", len(matchs))

	for _, v := range matchs {
		fmt.Printf("match %d. %s. Next match slot: %s, Next match: %v\n", v.Index, v.Round, v.NextMatchSlot, v.NextMatchIndex)
	}

	compatition.Matches = matchs
	compatition.NumOfRound = numOfRound
	compatition.NumOfMatch = len(matchs)
	if err := u.userrepository.InsertCompatition(compatition); err != nil {
		return err
	}

	return nil
}

func checkNumberPowerOfTwo(n int) int {
	return n & (n - 1)
}

// func (u *userUsecaseImpl) createRoundRobin(n int, compatitionID uint) []entities.Matches {
// 	lst := make([]int, n-1)
// 	matches := []entities.Matches{}
// 	for i := 0; i < len(lst); i++ {
// 		lst[i] = i + 2

// 	}
// 	if n%2 == 1 {
// 		lst = append(lst, 0) // 0 denotes a bye
// 		n++
// 	}
// 	for r := 1; r < n; r++ {
// 		fmt.Printf("Round %2d", r)
// 		lst2 := append([]int{1}, lst...)

// 		for i := 0; i < n/2; i++ {
// 			if lst[i] != 0 || lst2[n-1-i] != 0 {
// 				// fmt.Printf(" (%2d vs %-2d)", lst2[i], lst2[n-1-i])
// 				matches = append(matches, entities.Matches{
// 					CompetitionID: compatitionID,
// 					// DateTime    :,
// 					Team1ID     :,
// 					Team2ID     :,
// 					// Team1Goals  :,
// 					// Team2Goals  :,
// 					// Events      :,
// 					// GoalRecords :,
// 					// Result      :,
// 				})
// 				}
// 			}
// 		}
// 		fmt.Println()
// 		rotate(lst)
// 	}
// 	return nil
// }

func rotate(lst []int) {
	len := len(lst)
	last := lst[len-1]
	for i := len - 1; i >= 1; i-- {
		lst[i] = lst[i-1]
	}
	lst[0] = last
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
		if err.Error() == "record not found" {
			return nil, nil
		}
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
		if err.Error() == "record not found" {
			return &model.Team{}, nil
		}
		return &model.Team{}, err
	}
	memberList := []model.Member{}
	for _, member := range selectedTeams.TeamsMembers {
		memberList = append(memberList, model.Member{
			ID:            member.NormalUsers.ID,
			UsersID:       member.NormalUsers.UsersID,
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
			ID:           v.ID,
			Name:         v.Name,
			Format:       model.CompetitionFormat(v.Format),
			OrganizerID:  v.OrganizersID,
			StartDate:    v.StartDate,
			EndDate:      v.EndDate,
			AgeOver:      v.AgeOver,
			AgeUnder:     v.AgeUnder,
			Sex:          model.SexType(v.Sex),
			FieldSurface: model.FieldSurfaces(v.FieldSurface),
			Description:  v.Description,
			Status:       model.CompetitionStatus(v.Status),
			NumberOfTeam: v.NumberOfTeam,
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

// GetTeamList implements UserUsecase.
func (u *userUsecaseImpl) GetTeamsByOwnerID(in uint) ([]model.TeamList, error) {

	team := entities.Teams{
		// 0 is select all id
		OwnerID: in,
	}

	util.PrintObjInJson(team)

	teams, err := u.userrepository.GetTeams(&team, "id", false, -1, -1)
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

	util.PrintObjInJson(receiver)
	util.PrintObjInJson(team)

	for _, member := range team.TeamsMembers {
		fmt.Println(member.NormalUsers.UsersID)
		fmt.Println(receiver.UsersID)
		if member.NormalUsersID == receiver.ID {
			return errors.New("this user already in team")
		}
	}

	for _, requestSend := range team.RequestSends {
		if requestSend.ReceiverID == receiver.UsersID && requestSend.Status == "pending" {
			return errors.New("this user already invited")
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
		FirstNameThai: inUpdateModel.FirstNameThai,
		LastNameThai:  inUpdateModel.LastNameThai,
		FirstNameEng:  inUpdateModel.FirstNameEng,
		LastNameEng:   inUpdateModel.LastNameEng,
		Born:          inUpdateModel.Born,
		Height:        inUpdateModel.Height,
		Weight:        inUpdateModel.Weight,
		Sex:           inUpdateModel.Sex,
		Position:      inUpdateModel.Position,
		Nationality:   inUpdateModel.Nationality,
		Description:   inUpdateModel.Description,
		Phone:         inUpdateModel.Phone,
	}

	normalUser.ID = inNormalUserID

	util.PrintObjInJson(normalUser)

	if err := u.userrepository.UpdateNormalUser(normalUser); err != nil {
		return err
	}
	return nil
}

// CreateTeam implements UserUsecase.
func (u *userUsecaseImpl) CreateTeam(in *model.CreateTeam) error {

	normalUser, err := u.userrepository.GetNormalUser(&entities.NormalUsers{
		UsersID: in.OwnerID,
	})
	if err != nil {
		return err
	}

	//check required data
	requiredData := []string{}
	if normalUser.FirstNameThai == "" {
		requiredData = append(requiredData, "first_name_thai")
	}
	if normalUser.LastNameThai == "" {
		requiredData = append(requiredData, "last_name_thai")
	}
	if normalUser.FirstNameEng == "" {
		requiredData = append(requiredData, "first_name_eng")
	}
	if normalUser.LastNameEng == "" {
		requiredData = append(requiredData, "last_name_eng")
	}
	if normalUser.Born.IsZero() {
		requiredData = append(requiredData, "born")
	}
	if normalUser.Sex == "" {
		requiredData = append(requiredData, "sex")
	}
	if normalUser.Nationality == "" {
		requiredData = append(requiredData, "nationality")
	}
	if normalUser.Phone == "" {
		requiredData = append(requiredData, "phone")
	}
	if len(requiredData) != 0 {
		return &util.CreateTeamError{
			RequiredData: requiredData,
		}
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

func (u *userUsecaseImpl) Login(in *model.LoginUser) (string, model.LoginResponse, error) {

	if !isEmail(in.Email) {
		return "", model.LoginResponse{}, errors.New("invalid email format")
	}

	//get user from email
	user, err := u.userrepository.GetUserByEmail(in.Email)
	if err != nil {
		return "", model.LoginResponse{}, errors.New("incorrect email or password")
	}
	//compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password)); err != nil {
		return "", model.LoginResponse{}, errors.New("incorrect email or password")
	}

	loginResponse := model.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
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
			return "", model.LoginResponse{}, err
		}
		if normalUser.ImageProfilePath != "" {
			normalUser.ImageProfilePath = normalUser.ImageProfilePath[1:]
		}
		if normalUser.ImageCoverPath != "" {
			normalUser.ImageCoverPath = normalUser.ImageCoverPath[1:]
		}
		loginResponse.NormalUserID = normalUser.ID
		claims["normal_user_id"] = normalUser.ID
	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return "", model.LoginResponse{}, err
		}
		loginResponse.OrganizerID = organizer.ID
		claims["organizer_id"] = organizer.ID
	}
	//bypass
	jwtSecretKey := config.GetConfig().JwtSecretKey
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", model.LoginResponse{}, err
	}
	return t, loginResponse, nil
}

// GetUser implements UserUsecase.
func (u *userUsecaseImpl) GetUser(in uint) (model.User, error) {
	// get user from email
	user, err := u.userrepository.GetUserByID(in)
	if err != nil {
		return model.User{}, err
	}

	userModel := model.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,

		ImageProfilePath: user.ImageProfilePath,
		ImageCoverPath:   user.ImageCoverPath,
	}

	if user.Role == "normal" {
		normalUser, err := u.userrepository.GetNormalUserWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}

		userModel.NormalUserInfo = &model.NormalUserInfo{
			ID:            normalUser.ID,
			FirstNameThai: normalUser.FirstNameThai,
			LastNameThai:  normalUser.LastNameThai,
			FirstNameEng:  normalUser.FirstNameEng,
			LastNameEng:   normalUser.LastNameEng,
			Username:      normalUser.Username,
			Born:          normalUser.Born,
			Phone:         normalUser.Phone,
			Height:        normalUser.Height,
			Weight:        normalUser.Weight,
			Sex:           normalUser.Sex,
			Position:      normalUser.Position,
			Nationality:   normalUser.Nationality,
			Description:   normalUser.Description,
			Address: model.Address{
				HouseNumber: normalUser.Addresses.HouseNumber,
				Village:     normalUser.Addresses.Village,
				Subdistrict: normalUser.Addresses.Subdistrict,
				District:    normalUser.Addresses.District,
				PostalCode:  normalUser.Addresses.PostalCode,
				Country:     normalUser.Addresses.Country,
			},
		}
		// userModel.OrganizersInfo = model.OrganizersInfo{}

	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerWithAddressByUserID(user.ID)
		if err != nil {
			return model.User{}, err
		}
		userModel.OrganizersInfo = &model.OrganizersInfo{
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
		}
		// userModel.NormalUserInfo = model.NormalUserInfo{}

	}
	return userModel, nil
}
func (u *userUsecaseImpl) RegisterNormaluser(in *model.RegisterNormaluser) error {

	if !isEmail(in.Email) {
		return errors.New("invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}
	if isUsernameAlreadyInUser(in.Username, u.userrepository) {
		return errors.New("this username is already in use")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.Users{
		Email:            in.Email,
		Role:             "normal",
		Password:         string(hashedPassword),
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
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

func isUsernameAlreadyInUser(username string, u repositories.Userrepository) bool {
	if _, err := u.GetNormalUserByUsername(username); err != nil {
		return false
	}
	return true
}

func isPhoneAlreadyInUse(phone string, u repositories.Userrepository) bool {
	if _, err := u.GetNormalUser(&entities.NormalUsers{
		Phone: phone,
	}); err != nil {
		if _, err := u.GetOrganizer(&entities.Organizers{
			Phone: phone,
		}); err != nil {
			return false
		}
	}
	return true
}

func (u *userUsecaseImpl) RegisterOrganizer(in *model.RegisterOrganizer) error {

	if !isEmail(in.Email) {
		return errors.New("invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, u.userrepository) {
		return errors.New("this email is already in use")
	}

	if isPhoneAlreadyInUse(in.Phone, u.userrepository) {
		return errors.New("this phone is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.Users{
		Email:            in.Email,
		Role:             "organizer",
		Password:         string(hashedPassword),
		ImageProfilePath: "./images/default/defaultProfile.jpg",
		ImageCoverPath:   "./images/default/defaultCover.jpg",
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
			NormalUserInfo:   &model.NormalUserInfo{},
			ID:               e.ID,
			Email:            e.Email,
			Role:             e.Role,
			ImageProfilePath: e.ImageProfilePath,
			ImageCoverPath:   e.ImageCoverPath,
		}
		users_model = append(users_model, m)
	}

	return users_model, nil
}
