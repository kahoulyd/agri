package services

import (
	"Agri/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func (ps *ProductService) GetProducts(c *gin.Context) {
	rows, err := ps.db.Query("SELECT id, name, category, price FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}
	c.JSON(http.StatusOK, products)
}

func (ps *ProductService) AddProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO products (name, category, price) VALUES (?, ?, ?)"
	_, err := ps.db.Exec(query, product.Name, product.Category, product.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully", "product": product})
}

func (ps *ProductService) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE products SET name = ?, category = ?, price = ? WHERE id = ?"
	res, err := ps.db.Exec(query, product.Name, product.Category, product.Price, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (ps *ProductService) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM products WHERE id = ?"
	res, err := ps.db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
