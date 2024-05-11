package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
)

type UserService struct {
}

func (u *UserService) queryUserList(username, displayName string) []model.User {
	return nil
}

func (u *UserService) createUser(user model.User) {
	db := database.DB
	db.Create(&user)
}
