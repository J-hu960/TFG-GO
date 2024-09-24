package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	data "jordi.tfg.rewrite/internal"
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
	Email           string `json:"Email"`
	Password        PasswordS
	Profile_Picture string    `json:"profile_picture,omitempty"`
	Description     string    `json:"description,omitempty"`
	Role            string    `json:"role,omitempty"`
	Created_at      time.Time `json:"created_At"`
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

func (m UserModel) GetById(id int64) (*User, error) {
	var user User

	query := ` Select email,created_at,profile_pict,hashed_password,description,role
	  FROM users where pk_user = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&user.Email,
		&user.Created_at,
		&user.Profile_Picture,
		&user.Password.HashedPasswd,
		&user.Description,
		&user.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (m UserModel) DeleteById(id int64) error {
	query := `DELETE from users WHERE pk_user = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil

}

func (m UserModel) Update(user *User) error {
	query := `UPDATE users
	  set description = $1, role = $2
	  where pk_user = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, user.Description, user.Role, user.Pk_User)

	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetByMail(mail string) (*User, error) {
	var user User

	query := ` Select created_at,profile_pict,hashed_password,description,role
	  FROM users where email = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, mail).Scan(
		&user.Created_at,
		&user.Profile_Picture,
		&user.Password.HashedPasswd,
		&user.Description,
		&user.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (p *PasswordS) CreateHashedPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p.PlainText), 14)
	if err != nil {
		return err
	}

	p.HashedPasswd = string(bytes)

	return nil
}

func VerifyPassword(hashed, plaintex string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintex))
	return err == nil
}
