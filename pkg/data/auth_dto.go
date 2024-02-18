package data

type LoginWithEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,string"`
}

type LoginWithPhoneDTO struct {
	Phone string `json:"phone" validate:"required,e164"`
	Otp   string `json:"otp" validate:"required,string"`
}

type AccessTokenDTO struct {
	AccessToken string `json:"access_token"`
}
