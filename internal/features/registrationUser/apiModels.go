package registrationUser

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Authorization string `header:"Authorization" json:"-"`
}
