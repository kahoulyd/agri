package services

import (
	"Agri/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SupplierService struct {
	client *mongo.Client
}

func NewSupplierService(client *mongo.Client) *SupplierService {
	return &SupplierService{client: client}
}

func (ss *SupplierService) GetSuppliers(c *gin.Context) {
	collection := ss.client.Database("agriculture").Collection("suppliers")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var suppliers []models.Supplier
	for cursor.Next(context.Background()) {
		var s models.Supplier
		if err := cursor.Decode(&s); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		suppliers = append(suppliers, s)
	}
	c.JSON(http.StatusOK, suppliers)
}

func (ss *SupplierService) AddSupplier(c *gin.Context) {
	var supplier models.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := ss.client.Database("agriculture").Collection("suppliers")
	supplier.ID = primitive.NewObjectID().String() // Assign a new unique ID
	_, err := collection.InsertOne(context.TODO(), supplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func (ss *SupplierService) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	var updateData models.Supplier
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := ss.client.Database("agriculture").Collection("suppliers")
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updateData}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Supplier not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supplier updated successfully"})
}

func (ss *SupplierService) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	collection := ss.client.Database("agriculture").Collection("suppliers")
	filter := bson.M{"_id": objID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Supplier not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}
