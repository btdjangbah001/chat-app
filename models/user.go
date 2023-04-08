package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	About    string `json:"about"`
}

type UserLogin struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (user *User) CreateUser() (err error) {
	err = DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUser(id string) (*User, error) {
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

func GetUserByUsername() (*User, error) {
	var user User

	err := DB.Where("username = ?", user.Username).First(&user).Error
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