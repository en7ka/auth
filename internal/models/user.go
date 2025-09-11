package models

type UserInfoJwt struct {
	Username string `json:"username"`
	Role     bool   `json:"role"`
}
