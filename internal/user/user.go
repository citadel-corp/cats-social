package user

type User struct {
	ID             int64
	UID            string
	Email          string
	Name           string
	HashedPassword string
}
