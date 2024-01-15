package model

type User struct {
	UserID   uint   `gorm:"autoIncrement;primarykey" json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password`
	Role     Role   `json:"role"`
}
type UserRegisterRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password`
}
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLoginResponse struct {
	AccessToken string `json:"accessToken"`
	FullName    string `json:"fullName"`
	Role        Role   `json:"role"`
}
