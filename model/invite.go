package model

type Invite struct {
	GuildId     string `json:"guild_id"`
	IsPermanent bool   `json:"is_permanent"`
}
