package data

import (
	"context"
	"database/sql"
	"time"
)

type ReviewsModel struct {
	db *sql.DB
}

type Review struct {
	Pagek_relation int64     `json:"pk_review"`
	Id_user        int64     `json:"pk_creator"`
	Id_project     int64     `json:"pk_project"`
	Created_at     time.Time `json:"created_at"`
	Content        string    `json:"content"`
}

func (m ReviewsModel) Insert(review Review) error {

	query := ` INSERT INTO reviews (id_user,id_project,content)
	values($1,$2,$3)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, review.Id_user, review.Id_project, review.Content)
	return err

}

func (m ReviewsModel) Delete(id int64) error {

	query := ` DELETE FROM reviews WHERE pk_relation = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err

}
