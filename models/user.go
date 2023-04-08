package models

import "github.com/gorilla/websocket"

type User struct {
	ID       uint            `json:"id" gorm:"primary_key"`
	Username string          `json:"username" gorm:"unique_index;not null"`
	Email    string          `json:"email" gorm:"unique_index;not null"`
	Password string          `json:"password" gorm:"not null"`
	About    string          `json:"about"`
	Ws       *websocket.Conn `json:"ws"`
}

type UserSignUp struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserLogin struct {
	UserField string `json:"user_field"`
	Password  string `json:"password"`
}

func (user *User) CreateUser() (err error) {
	err = DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUser(id uint) (*User, error) {
	var user User

	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user *User) UpdateUser(updateUser *User) (err error) {

	err = DB.Model(&user).Where("id = ?", user.ID).Updates(updateUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) DeleteUser() (err error) {
	err = DB.Where("id = ?", user.ID).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User

	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User

	err := DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersForGroup(groupID uint) (*[]User, error) {
	var users []User
	err := DB.Table("users").Joins("JOIN memberships ON memberships.user_id = users.id").Where("memberships.group_id = ?", groupID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func UsernameOrEmailExists(user *UserSignUp) (bool, error) {
	err := DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateUserWebsocket(user *User, ws *websocket.Conn) error {
	err := DB.Model(&user).Where("id = ?", user.ID).Update("ws", ws).Error
	if err != nil {
		return err
	}
	return nil
}