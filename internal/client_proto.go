package main

import (
	"Agri/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewSupplierServiceProtoClient(conn)

	// Add a Supplier
	addResp, err := client.P_AddSupplier(context.Background(), &proto.AddSupplierRequest{
		Name:    "John's Farm",
		Contact: "john@example.com",
	})
	if err != nil {
		log.Fatalf("Error adding supplier: %v", err)
	}
	fmt.Println("Added Supplier:", addResp)

	// Get all Suppliers
	time.Sleep(1 * time.Second)
	suppliers, err := client.P_GetSuppliers(context.Background(), &proto.SupplierRequest{})
	if err != nil {
		log.Fatalf("Error getting suppliers: %v", err)
	}
	fmt.Println("Supplier List:", suppliers)

	// Update Supplier
	updateResp, err := client.P_UpdateSupplier(context.Background(), &proto.UpdateSupplierRequest{
		Id:      addResp.Id,
		Name:    "John's Organic Farm",
		Contact: "johnorganic@example.com",
	})
	if err != nil {
		log.Fatalf("Error updating supplier: %v", err)
	}
	fmt.Println("Updated Supplier:", updateResp)

	// Delete Supplier
	deleteResp, err := client.P_DeleteSupplier(context.Background(), &proto.DeleteSupplierRequest{Id: addResp.Id})
	if err != nil {
		log.Fatalf("Error deleting supplier: %v", err)
	}
	fmt.Println(deleteResp.Message)
}
