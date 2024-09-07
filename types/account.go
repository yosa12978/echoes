package types

type Account struct {
	Id       string
	Username string
	Password string
	Salt     string
	Created  string
	IsAdmin  bool
}
