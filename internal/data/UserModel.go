package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	db *sql.DB
}
type PasswordS struct {
	PlainText    string `json:"-"`
	HashedPasswd string `json:"-"`
}

type User struct {
	Pk_User         int64  `json:"-"`
	Phone           string `json:"Phone"`
	Email           string `json:"Email"`
	Password        PasswordS
	Profile_Picture string `json:"profile_picture"`
	Description     string `json:"description"`
	Role            string `json:"role"`
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (email,hashed_password) VALUES($1,$2)`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	fmt.Print(user.Password.HashedPasswd)

	_, err := m.db.ExecContext(ctx, query, user.Email, user.Password.HashedPasswd)
	if err != nil {
		return err
	}
	return nil
}

func (p *PasswordS) CreateHashedPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.PlainText), 14)
	if err != nil {
		return err
	}

	p.HashedPasswd = string(bytes)

	return nil

}

func (p *PasswordS) VerifyPassword() bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.HashedPasswd), []byte(p.PlainText))
	return err == nil
}
