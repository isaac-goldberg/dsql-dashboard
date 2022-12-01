package typings

type UserData struct {
	UserId        string
	Username      string
	Discriminator string
	IconURL       string
}

type TemplatingData struct {
	Title string
	User  UserData
}

type RouteLocalsKeys_User struct{}
