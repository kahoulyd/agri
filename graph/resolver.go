package graph

import (
	"Agri/models"
	"context"
	"database/sql"
	"fmt"
)

// Resolver struct
type Resolver struct {
	DB *sql.DB
}

// Get all products
func (r *Resolver) GetProducts(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.DB.Query("SELECT id, name, category, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

// Add a new product
func (r *Resolver) AddProduct(ctx context.Context, name string, category string, price float64) (*models.Product, error) {
	query := "INSERT INTO products (name, category, price) VALUES (?, ?, ?)"
	res, err := r.DB.Exec(query, name, category, price)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	product := &models.Product{
		ID:       int(id),
		Name:     name,
		Category: category,
		Price:    price,
	}
	return product, nil
}

// Update a product
func (r *Resolver) UpdateProduct(ctx context.Context, id int, name *string, category *string, price *float64) (*models.Product, error) {
	query := "UPDATE products SET name = ?, category = ?, price = ? WHERE id = ?"
	_, err := r.DB.Exec(query, name, category, price, id)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:       id,
		Name:     *name,
		Category: *category,
		Price:    *price,
	}, nil
}

// Delete a product
func (r *Resolver) DeleteProduct(ctx context.Context, id int) (string, error) {
	query := "DELETE FROM products WHERE id = ?"
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return "", err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return "", fmt.Errorf("Product not found")
	}

	return "Product deleted successfully", nil
}
