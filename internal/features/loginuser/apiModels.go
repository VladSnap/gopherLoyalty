package loginuser

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Authorization string `header:"Authorization" json:"-"`
}
