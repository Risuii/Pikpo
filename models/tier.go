package models

type Tier struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	UserId int    `json:"user_id"`
	Tier   string `json:"tier"`
}
