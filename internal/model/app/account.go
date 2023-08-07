package app

type Account struct {
	Id           int
	Username     string
	HashPassword string
	Role         Role
}
