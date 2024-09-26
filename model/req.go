package model

type UserLogin struct {
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	VerifyKey   string `json:"verify_key"`
	VerifyValue string `json:"verify_value"`
}
