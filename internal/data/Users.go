package data

import (
	"context"
	"database/sql"
	"errors"
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
	Email           string `json:"Email"`
	Password        PasswordS
	Profile_Picture string    `json:"profile_picture,omitempty"`
	Description     string    `json:"description,omitempty"`
	Role            string    `json:"role,omitempty"`
	Created_at      time.Time `json:"created_At"`
	Version         int       `json:"version"`
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
			return nil, ErrNotFound
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

	query := ` Select pk_user,created_at,profile_pict,hashed_password,description,role
	  FROM users where email = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, mail).Scan(
		&user.Pk_User,
		&user.Created_at,
		&user.Profile_Picture,
		&user.Password.HashedPasswd,
		&user.Description,
		&user.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
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

func (m UserModel) HasAlreadyLikedProject(userId, projectId int64) (bool, error) {
	query := "SELECT pk_relation FROM user_liked_projects where id_user = $1 AND id_project = $2"

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var pkRelation int64

	err := m.db.QueryRowContext(ctx, query, userId, projectId).Scan(&pkRelation)

	print("Error: ", err)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil

}

func (m UserModel) HasAlreadyDislikedProject(userId, projectId int64) (bool, error) {
	query := "SELECT pk_relation FROM user_disliked_projects where id_user = $1 AND id_project = $2"

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var pkRelation int64

	err := m.db.QueryRowContext(ctx, query, userId, projectId).Scan(&pkRelation)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil

}

func (m UserModel) CreateUserProjectLikeRelation(userId, projectId int64) error {
	query := `INSERT INTO user_liked_projects (id_user,id_project) VALUES($1,$2)`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, userId, projectId)
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) DeleteUserProjectLikeRelation(iduser, idproject int64) error {
	query := `DELETE FROM user_liked_projects WHERE id_user=$1 AND id_project = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, iduser, idproject)
	if err != nil {
		return err
	}
	return nil
}
func (m UserModel) CreateUserProjectDisLikeRelation(userId, projectId int64) error {
	query := `INSERT INTO user_disliked_projects (id_user,id_project) VALUES($1,$2)`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, userId, projectId)
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) DeleteUserProjectDislikeRelation(iduser, idproject int64) error {
	query := `DELETE FROM user_disliked_projects WHERE id_user=$1 AND id_project = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, iduser, idproject)
	if err != nil {
		return err
	}
	return nil
}
