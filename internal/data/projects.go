package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type ProjectModel struct {
	db *sql.DB
}

type Project struct {
	Pk_project     int64     `json:"pk_project"`
	Name           string    `json:"project_name"`
	CreatedAt      time.Time `json:"created_at"`
	Photos         []string  `json:"photos"`
	Link_web       string    `json:"link_web"`
	Description    string    `json:"description"`
	Likes          int64     `json:"likes"`
	Disikes        int64     `json:"dislikes"`
	FoundsRecieved int64     `json:"founds_recieved"`
	FoundsExpected int64     `json:"founds_expected"`
	Category       []string  `json:"categories"`
	IdCreator      int64     `json:"id_creator"`
	Version        int       `json:"version"`
}

func (m ProjectModel) Insert(project Project) error {
	query := `INSERT INTO projects (name,photos,link_web,description,founds_expected,category,id_creator)
	values($1,$2,$3,$4,$5,$6,$7)
 	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	args := []any{project.Name, pq.Array(project.Photos), project.Link_web, project.Description, project.FoundsExpected, pq.Array(project.Category), project.IdCreator}

	_, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (m ProjectModel) GetById(id int64) (*Project, error) {
	var project Project
	query := `SELECT pk_project, name, created_at,photos,link_web,description,likes,dislikes,founds_recieved,founds_expected,category,id_creator
	from projects where pk_project = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&project.Pk_project,
		&project.Name,
		&project.CreatedAt,
		pq.Array(&project.Photos),
		&project.Link_web,
		&project.Description,
		&project.Likes,
		&project.Disikes,
		&project.FoundsRecieved,
		&project.FoundsRecieved,
		pq.Array(&project.Category),
		&project.IdCreator,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &project, nil
}

func (m ProjectModel) UpdateProject(project *Project) error {
	query := ` UPDATE projects set
	   name = $1, photos=$2, link_web=$3,description=$4, founds_expected=$5,category=$6
	   where pk_project = $7
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, project.Name, pq.Array(project.Photos), project.Link_web, project.Description, project.FoundsExpected, pq.Array(project.Category), project.Pk_project)

	return err
}

func (m ProjectModel) DeleteProject(id int64) error {
	query := `DELETE FROM projects WHERE pk_project = $q
	`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m ProjectModel) GetAll(name string, categories []string, filters Filters) ([]Project, error) {
	var projects []Project
	offset := (filters.Page - 1) * filters.PageSize

	query := fmt.Sprintf(`
	SELECT pk_project, name, created_at, photos, link_web, description, likes, dislikes, founds_recieved, founds_expected, category, id_creator
	FROM projects 
	WHERE (LOWER(name)=LOWER($1) OR $1 = '')
    AND (category @> $2 OR category ='{}')
	ORDER BY %s %s, pk_project ASC
	LIMIT $3 OFFSET $4
	`, filters.SortColumn(), filters.SortOrder())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, name, pq.Array(categories), filters.PageSize, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.Pk_project,
			&project.Name,
			&project.CreatedAt,
			pq.Array(&project.Photos),
			&project.Link_web,
			&project.Description,
			&project.Likes,
			&project.Disikes,
			&project.FoundsRecieved,
			&project.FoundsExpected,
			pq.Array(&project.Category),
			&project.IdCreator,
		)

		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (m ProjectModel) AddLike(id int64) error {
	query := `UPDATE projects SET likes = likes +1 
	  WHERE pk_project = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err

}

func (m ProjectModel) AddDislike(id int64) error {
	query := `UPDATE projects SET dislikes = dislikes +1 
	  WHERE pk_project = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err

}

func (m ProjectModel) SubLike(id int64) error {
	query := `UPDATE projects SET likes = likes - 1 
	  WHERE pk_project = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err

}

func (m ProjectModel) SubDislike(id int64) error {
	query := `UPDATE projects SET dislikes = dislikes - 1 
	  WHERE pk_project = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	return err

}
