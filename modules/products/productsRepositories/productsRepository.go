package productsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/Kamila3820/go-shop-tutorial/config"
	"github.com/Kamila3820/go-shop-tutorial/modules/entities"
	"github.com/Kamila3820/go-shop-tutorial/modules/files/filesUsecases"
	"github.com/Kamila3820/go-shop-tutorial/modules/products"
	"github.com/Kamila3820/go-shop-tutorial/modules/products/productsPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) ([]*products.Product, int)
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
}

type productsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "product_category" "pc" ON "pc"."category_id" = "c"."id"
				WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "image" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "product" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";
	`
	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get product failed %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed %v", err)
	}

	return product, nil
}

func (r *productsRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productsPatterns.FindProductBuilder(r.db, req)
	engineer := productsPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	return result, count
}

func (r *productsRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.InsertProductBuilder(r.db, req)
	productId, err := productsPatterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productsRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productsPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
