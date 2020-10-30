package chat

type UserRepository interface {
	GetUser(channel string, user string) ChatUser
	SaveUser(channel string, cu ChatUser)
}
