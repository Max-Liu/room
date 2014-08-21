package room

type User struct {
	Name string
	Cmd
	Msg Message
}

func NewUser() *User {
	return &User{}
}
