package data

import "database/sql"

type Models struct {
	Users    UserModel
	Projects ProjectModel
	Reviews  ReviewsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{db: db},
		Projects: ProjectModel{db: db},
		Reviews:  ReviewsModel{db: db},
	}
}
