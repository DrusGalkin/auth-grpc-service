package models

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	HashPassword []byte `json:"password_hash"`
	Role         string `json:"role"`
}
