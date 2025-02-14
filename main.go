package main

import (
	"Agri/graph"
	_ "Agri/models"
	"Agri/services"
	"context"
	"database/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	_ "net/http"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=agriculture sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	productService := services.NewProductService(db)
	supplierService := services.NewSupplierService(mongoClient)

	r := gin.Default()

	//products
	r.GET("/products", productService.GetProducts)
	r.POST("/products", productService.AddProduct)
	r.PUT("/products/:id", productService.UpdateProduct)
	r.DELETE("/products/:id", productService.DeleteProduct)

	//suppliers
	r.GET("/suppliers", supplierService.GetSuppliers)
	r.POST("/suppliers", supplierService.AddSupplier)
	r.PUT("/suppliers/:id", supplierService.UpdateSupplier)
	r.DELETE("/suppliers/:id", supplierService.DeleteSupplier)

	//Graphql
	resolver := graph.Resolver{DB: db}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	// Routes
	r.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":8081")
}
