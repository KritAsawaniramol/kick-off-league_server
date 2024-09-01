package authUsecase

import (
	"errors"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"kickoff-league.com/config"
	"kickoff-league.com/entities"
	model "kickoff-league.com/models"
	"kickoff-league.com/repositories"
	"kickoff-league.com/util"
)

type authUsecaseImpl struct {
	repository repositories.Repository
}

func NewAuthUsecaseImpl(
	userrepository repositories.Repository,
) AuthUsecase {
	return &authUsecaseImpl{
		repository: userrepository,
	}
}

// RegisterOrganizer implements AuthUsecase.
func (a *authUsecaseImpl) RegisterOrganizer(in *model.RegisterOrganizer) error {

	if !isEmail(in.Email) {
		return errors.New("error: invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, a.repository) {
		return errors.New("error: this email is already in use")
	}

	if isPhoneAlreadyInUse(in.Phone, a.repository) {
		return errors.New("error: this phone is already in use")
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

	if err := a.repository.InsertUserWihtOrganizerAndAddress(organizer, user); err != nil {
		return err
	}

	return nil
}

// Login implements AuthUsecase.
func (a *authUsecaseImpl) Login(in *model.LoginUser) (string, model.LoginResponse, error) {
	if !util.IsEmail(in.Email) {
		return "", model.LoginResponse{}, errors.New("error: invalid email format")
	}

	//get user from email
	user, err := a.repository.GetUserByEmail(in.Email)
	if err != nil {
		return "", model.LoginResponse{}, errors.New("error: user not found")
	}
	//compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password)); err != nil {
		return "", model.LoginResponse{}, errors.New("error: incorrect email or password")
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
		normalUser, err := a.repository.GetNormalUserDetails(
			&entities.NormalUsers{
				UsersID: user.ID,
			},
		)
		if err != nil {
			return "", model.LoginResponse{}, err
		}
		if user.ImageProfilePath != "" {
			user.ImageProfilePath = user.ImageProfilePath[1:]
		}
		if user.ImageCoverPath != "" {
			user.ImageCoverPath = user.ImageCoverPath[1:]
		}
		loginResponse.NormalUserID = normalUser.ID
		claims["normal_user_id"] = normalUser.ID
	} else if user.Role == "organizer" {
		organizer, err := a.repository.GetOrganizer(
			&entities.Organizers{
				UsersID: user.ID,
			},
		)
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

// RegisterNormaluser implements AuthUsecase.
func (a *authUsecaseImpl) RegisterNormaluser(in *model.RegisterNormaluser) error {

	if !isEmail(in.Email) {
		return errors.New("error: invalid email format")
	}

	if isEmailAlreadyInUse(in.Email, a.repository) {
		return errors.New("error: this email is already in use")
	}
	if isUsernameAlreadyInUser(in.Username, a.repository) {
		return errors.New("error: this username is already in use")
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

	if err := a.repository.InsertUserWihtNormalUserAndAddress(normalUser, user); err != nil {
		return err
	}

	return nil
}

func isEmailAlreadyInUse(email string, u repositories.Repository) bool {
	if _, err := u.GetUserByEmail(email); err != nil {
		return false
	}
	return true
}

func isUsernameAlreadyInUser(username string, u repositories.Repository) bool {
	if _, err := u.GetNormalUserByUsername(username); err != nil {
		return false
	}
	return true
}

func isPhoneAlreadyInUse(phone string, u repositories.Repository) bool {
	if _, err := u.GetNormalUserDetails(&entities.NormalUsers{
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

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
