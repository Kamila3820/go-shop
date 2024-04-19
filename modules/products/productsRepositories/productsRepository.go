package productsRepositories

import (
	"github.com/Kamila3820/go-shop-tutorial/config"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
}

type productsRepository struct {
	db  *sqlx.DB
	cfg config.IConfig
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig) IProductsRepository {
	return &productsRepository{
		db:  db,
		cfg: cfg,
	}
}
