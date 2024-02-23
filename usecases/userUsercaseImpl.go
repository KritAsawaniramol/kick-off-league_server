package usecases

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"kickoff-league.com/config"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
)

type userUsecaseImpl struct {
	userrepository repositories.Userrepository
}

// UpdateNormalUser implements UserUsecase.
func (u *userUsecaseImpl) UpdateNormalUser(inUpdateModel *model.UpdateNormalUser, inUserID uint) error {

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
	if err := u.userrepository.UpdateNormalUser(normalUser, inUserID); err != nil {
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

	team := entities.Team{
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
func (u *userUsecaseImpl) Register(in *model.RegisterUser) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	in.Password = string(hashedPassword)

	user := &entities.User{
		Email:    in.Email,
		Role:     in.Role,
		Password: in.Password,
	}

	if in.Role == "normal" {
		normalUser := &entities.NormalUser{
			UserID: user.ID,
		}
		if err := u.userrepository.InsertUserWihtNormalUserAndAddress(normalUser, user); err != nil {
			return err
		}
	} else if in.Role == "organizer" {
		organizer := &entities.Organizer{
			UserID: user.ID,
		}
		if err := u.userrepository.InsertUserWihtOrganizerAndAddress(organizer, user); err != nil {
			return err
		}
	}

	return nil
}

// GetUsers implements UserUsecase.
func (u *userUsecaseImpl) GetUsers() ([]entities.User, error) {
	users := []entities.User{}
	users, err := u.userrepository.GetUsers()
	if err != nil {
		return []entities.User{}, err
	}
	return users, nil
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
