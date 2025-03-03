package models

type RefreshToken struct {
	ID    uint   `json:"ID"`
	Token string `json:"token"`
}
