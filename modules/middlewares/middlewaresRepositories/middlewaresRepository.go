package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewaresRepository interface {
	FindAccessToken(userId string, accessToken string) bool
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId string, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;
	`
	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}
