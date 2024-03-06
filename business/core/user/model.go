package user

type NewUser struct {
	Username string
	Password string
	Email    string
}

type NewToken struct {
	Username string
	Password string
}
