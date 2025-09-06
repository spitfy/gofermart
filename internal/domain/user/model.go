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
