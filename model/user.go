package model

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
