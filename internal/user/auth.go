package user

import (
	"fmt"
	hdr "marketplace/internal/handlers"
	"marketplace/internal/utils"
)

// SignIn авторизует уже зарегистрированного пользователя
func (repo *UserDBRepository) SignIn(usr *User) (*User, int, error) {
	exists, err := utils.CheckUser(repo.dtb, usr.Username)
	if err != nil {
		return nil, hdr.InternalServerErrorCode, err
	}

	if !exists {
		return nil, hdr.UnauthorizedCode, nil
	}

	passwordHash, err := GetPasswordHash(repo.dtb, usr.Username)
	if err != nil {
		return nil, hdr.InternalServerErrorCode, err
	}
	check, err := CheckPassword(usr.Password, passwordHash)
	if err != nil {
		return nil, hdr.InternalServerErrorCode, err
	}
	if !check {
		return nil, hdr.UnauthorizedCode, fmt.Errorf("password is incorrect")
	}

	thisUser := User{Username: usr.Username, Password: passwordHash}
	return &thisUser, hdr.OKCode, nil
}

// SignUp регистрирует нового пользователя
func (repo *UserDBRepository) SignUp(usr *User) (*User, int, error) {
	exists, err := utils.CheckUser(repo.dtb, usr.Username)
	if err != nil {
		return nil, hdr.InternalServerErrorCode, err
	}

	if exists {
		return nil, hdr.BadRequestCode, nil
	}

	thisUser, err := CreateUser(repo.dtb, usr)
	if err != nil {
		return nil, hdr.InternalServerErrorCode, err
	}

	return thisUser, hdr.OKCode, nil
}
