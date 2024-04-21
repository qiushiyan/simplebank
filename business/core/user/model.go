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

type NewVerifyEmail struct {
	Email string
}

type FinishVerifyEmail struct {
	Id   int64
	Code string
}
