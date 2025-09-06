package user

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID       int
	Login    string
	Password string
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
