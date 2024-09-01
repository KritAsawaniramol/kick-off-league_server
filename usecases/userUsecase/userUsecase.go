package userUsecase

import (
	model "kickoff-league.com/models"
)

type UserUsecase interface {
	GetUsers() ([]model.User, error)
	GetUser(in uint) (model.User, error)
	RemoveImageProfile(userID uint) error
	RemoveImageCover(userID uint) error
	UpdateImageCover(userID uint, newImagePath string) error
	UpdateImageProfile(userID uint, newImagePath string) error
}
