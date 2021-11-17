package dto

type User struct {
	Uid           int    `boil:"uid" json:"uid"`
	LoginName     string `boil:"login_name" json:"login_name,omitempty"`
	Nickname      string `boil:"nickname" json:"nickname"`
	ShowLoginName bool   `boil:"show_login_name" json:"show_login_name"`
}

type Contact struct {
	Uid         int    `boil:"uid" json:"uid"`
	DisplayName string `boil:"display_name" json:"display_name"`
	Cid         int    `boil:"cid" json:"cid"`
}

type Chats = []int

type Message struct {
	Mid     int    `boil:"mid" json:"mid"`
	Cid     int    `boil:"cid" json:"cid"`
	Uid     int    `boil:"uid" json:"uid"`
	Mime    string `boil:"mime" json:"mime"`
	Content string `boil:"content" json:"content"`
}
