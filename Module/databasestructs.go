package Module

type AdminInfo struct {
	Admin_id  int
	Admin_uid string
}

type UserInfo struct {
	Uid        int
	TelegramId string
	Age        string
	Role       string
	Height     string
	Bodytype   string
	Size       string
}

type WelcomeMessage struct {
	Group_username string
	Group_welcome  string
	Ask_role       int
}
