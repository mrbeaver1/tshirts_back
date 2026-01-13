package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/mrbeaver1/tshirts_back/internal/domain/entity"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) FindAll() []*entity.Product {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, name, color_id, size_id, quantity, available, parent_id, created_at, updated_at, deleted_at
		FROM products
		WHERE deleted_at IS NULL
	ORDER BY id
	`)
	if err != nil {
		fmt.Printf("Error querying products: %v\n", err)
		return nil
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		var parentId *int64
		var createdAt *time.Time
		var updatedAt *time.Time
		var deletedAt *time.Time

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.ColorId,
			&product.SizeId,
			&product.Quantity,
			&product.Available,
			&parentId,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			fmt.Printf("Error scanning product row: %v\n", err)
			continue
		}

		if parentId != nil {
			pid := uint64(*parentId)
			product.ParentId = &pid
		} else {
			product.ParentId = nil
		}

		product.CreatedAt = createdAt
		product.UpdatedAt = updatedAt
		product.DeletedAt = deletedAt

		products = append(products, product)
	}

	return products
}

func (r *PostgresProductRepository) FindOneById(id uint64) *entity.Product {
	row := r.db.QueryRow(context.Background(), `
		SELECT id, name, color_id, size_id, quantity, available, parent_id, created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	product := &entity.Product{}
	var parentId *int64
	var createdAt *time.Time
	var updatedAt *time.Time
	var deletedAt *time.Time

	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.ColorId,
		&product.SizeId,
		&product.Quantity,
		&product.Available,
		&parentId,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		fmt.Printf("Error scanning product: %v\n", err)
		return nil
	}

	if parentId != nil {
		pid := uint64(*parentId)
		product.ParentId = &pid
	} else {
		product.ParentId = nil
	}

	product.CreatedAt = createdAt
	product.UpdatedAt = updatedAt
	product.DeletedAt = deletedAt

	return product
}

func (r *PostgresProductRepository) Create(p *entity.Product) {
	var parentId *int64
	if p.ParentId != nil {
		pid := int64(*p.ParentId)
		parentId = &pid
	}

	query := `
		INSERT INTO products (name, color_id, size_id, quantity, available, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		p.Name,
		p.ColorId,
		p.SizeId,
		p.Quantity,
		p.Available,
		parentId,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&p.Id)
	if err != nil {
		fmt.Printf("Error inserting product: %v\n", err)
	}
}

func (r *PostgresProductRepository) Update(p *entity.Product) {
	var parentId *int64
	if p.ParentId != nil {
		pid := int64(*p.ParentId)
		parentId = &pid
	}

	query := `
		UPDATE products
	SET name = $2, color_id = $3, size_id = $4, quantity = $5, available = $6, parent_id = $7, updated_at = $8
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		p.Id,
		p.Name,
		p.ColorId,
		p.SizeId,
		p.Quantity,
		p.Available,
		parentId,
		time.Now().UTC(),
	)
	if err != nil {
		fmt.Printf("Error updating product: %v\n", err)
	}
}

func (r *PostgresProductRepository) Delete(id uint64) {
	query := `
	UPDATE products
	SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Exec(context.Background(), query, time.Now().UTC(), id)
	if err != nil {
		fmt.Printf("Error deleting product: %v\n", err)
	}
}
