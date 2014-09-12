package room

type User struct {
	Name string
	Cmd
	Msg Message
}
type Message struct {
	Content string
}

func NewUser() *User {
	return &User{}
}
