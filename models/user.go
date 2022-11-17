package models

import "time"

type User struct {
	ID         int    `json:"id"`
	UsernameT3 string `json:"username_T3"`
	UsernameT2 string `json:"username_T2"`
	UsernameT1 string `json:"username_T1"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	ExpiredAt  time.Time
	Tier       []Tier `foreignKey:"UserId"`
}
