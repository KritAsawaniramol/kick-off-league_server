package usecases

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"kickoff-league.com/config"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
)

type userUsecaseImpl struct {
	userrepository repositories.Userrepository
}

// IgnoreAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) IgnoreAddMemberRequest(inReqID uint, userID uint) error {
	addMemberRequest, err := u.userrepository.GetAddMemberRequestByID(inReqID)
	if err != nil {
		return err
	}

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

	addMemberRequest, err := u.userrepository.GetAddMemberRequestByID(inReqID)
	if err != nil {
		return err
	}

	if addMemberRequest.ReceiverID != userID {
		log.Printf("ReceiverID: %d, userID: %d\n", addMemberRequest.ReceiverID, userID)
		return errors.New("UserID does not match ReceiverID")
	}

	normalUser, err := u.userrepository.GetNormalUserByUserID(userID)
	if err != nil {
		return err
	}

	team, err := u.userrepository.GetTeamByID(addMemberRequest.TeamsID)
	if err != nil {
		return err
	}

	if err := u.userrepository.UpdateAddMemberRequestStatusByID(inReqID, "accepted"); err != nil {
		return err
	}

	normalUser.Teams = append(normalUser.Teams, *team)

	// normalUser := &entities.NormalUser{
	// 	TeamID:       addMemberRequest.TeamsID,
	// 	NormalUserID: addMemberRequest.ReceiverID,
	// }

	if err := u.userrepository.UpdateNormalUser(normalUser); err != nil {
		return err
	}

	return nil
}

// CreateAddMemberRequest implements UserUsecase.
func (u *userUsecaseImpl) SendAddMemberRequest(inAddMemberRequest *model.AddMemberRequest, inUserID uint) error {
	// Get Receiver normaluser
	receiver, err := u.userrepository.GetNormalUserByUsername(inAddMemberRequest.ReceiverUsername)
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

	for _, member := range team.Member {
		if member.UserID == receiver.UserID {
			return errors.New("this request receiver is already in team")
		}
	}

	for _, requestSend := range team.RequestSend {
		if requestSend.ReceiverID == receiver.UserID && requestSend.Status == "pending" {
			return errors.New("team have already sent a request to this receiver")
		}
	}

	addMemberRequest := &entities.AddMemberRequest{
		TeamsID:    inAddMemberRequest.TeamID,
		ReceiverID: receiver.UserID,
		Status:     "pending",
	}

	// Create Request
	if err := u.userrepository.InsertAddMemberRequest(addMemberRequest); err != nil {
		return err
	}

	return nil

}

// UpdateNormalUser implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inNormalUserID uint) error {

	// normalUser, err := u.userrepository.GetNormalUserByUserID(inUserID)
	// if err != nil {
	// 	return err
	// }
	normalUser := &entities.NormalUser{
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

	if err := u.userrepository.UpdateNormalUser(normalUser); err != nil {
		return err
	}
	return nil
}

// CreateTeam implements UserUsecase.
func (u *userUsecaseImpl) CreateTeam(in *model.CreaetTeam) error {

	normalUser, err := u.userrepository.GetNormalUserByUserID(in.OwnerID)
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
		Name:        in.Name,
		OwnerID:     in.OwnerID,
		Member:      []entities.NormalUser{*normalUser},
		Compatition: []entities.Compatition{},
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

// GetUserByPhone implements UserUsecase.
// func (u *userUsecaseImpl) GetUserByPhone(in string) (model.User, error) {

// 	user_entities, err := u.userrepository.GetNormalUserByPhone(in)
// 	if err != nil {
// 		return model.NormalUser{}, err
// 	}

// 	user_model := model.NormalUser{
// 		UserID:        user_entities.UserID,
// 		FirstNameThai: user_entities.FirstNameThai,
// 		LastNameThai:  user_entities.LastNameThai,
// 		FirstNameEng:  user_entities.FirstNameEng,
// 		LastNameEng:   user_entities.LastNameEng,
// 		Born:          user_entities.Born,
// 		Height:        user_entities.Height,
// 		Weight:        user_entities.Weight,
// 		Sex:           user_entities.Sex,
// 		Position:      user_entities.Position,
// 		Nationality:   user_entities.Nationality,
// 		Phone:         user_entities.Phone,
// 		AddressID:     user_entities.AddressID,
// 	}
// 	return user_model, nil
// }

func NewUserUsercaseImpl(
	userrepository repositories.Userrepository,
) UserUsecase {
	return &userUsecaseImpl{
		userrepository: userrepository,
	}
}

func (u *userUsecaseImpl) Login(in *model.LoginUser) (string, model.User, error) {
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
		normalUser, err := u.userrepository.GetNormalUserByUserID(user.ID)
		if err != nil {
			return "", model.User{}, err
		}
		userModel.Datail = map[string]interface{}{
			"normal": normalUser,
		}
		claims["normalUser_id"] = userModel.ID
	} else if user.Role == "organizer" {
		organizer, err := u.userrepository.GetOrganizerByUserID(user.ID)
		if err != nil {
			return "", model.User{}, err
		}
		userModel.Datail = map[string]interface{}{
			"organizer": organizer,
		}
		claims["organizer_id"] = userModel.ID
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	in.Password = string(hashedPassword)

	user := &entities.User{
		Email:    in.Email,
		Role:     "normal",
		Password: in.Password,
	}

	normalUser := &entities.NormalUser{
		UserID:   user.ID,
		Username: in.Username,
	}
	if err := u.userrepository.InsertUserWihtNormalUserAndAddress(normalUser, user); err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) RegisterOrganizer(in *model.RegisterUser) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	in.Password = string(hashedPassword)

	user := &entities.User{
		Email:    in.Email,
		Role:     "organizer",
		Password: in.Password,
	}

	organizer := &entities.Organizer{
		UserID: user.ID,
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
