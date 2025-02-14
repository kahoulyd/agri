package proto

import (
	"context"
	"errors"
	"log"
	"sync"
)

// In-memory data store
var suppliers = make(map[string]*Supplier)
var mu sync.Mutex

// SupplierService struct
type SupplierServiceProto struct {
	UnimplementedSupplierServiceProtoServer
}

// Get all suppliers
func (s *SupplierServiceProto) P_GetSuppliers(ctx context.Context, req *SupplierRequest) (*SupplierList, error) {
	mu.Lock()
	defer mu.Unlock()

	var supplierList []*Supplier
	for _, supplier := range suppliers {
		supplierList = append(supplierList, supplier)
	}

	return &SupplierList{Suppliers: supplierList}, nil
}

// Add a supplier
func (s *SupplierServiceProto) P_AddSupplier(ctx context.Context, req *AddSupplierRequest) (*Supplier, error) {
	mu.Lock()
	defer mu.Unlock()

	id := req.Name + "_id"
	supplier := &Supplier{
		Id:      id,
		Name:    req.Name,
		Contact: req.Contact,
	}
	suppliers[id] = supplier

	log.Println("Supplier added:", supplier)
	return supplier, nil
}

// Update a supplier
func (s *SupplierServiceProto) P_UpdateSupplier(ctx context.Context, req *UpdateSupplierRequest) (*Supplier, error) {
	mu.Lock()
	defer mu.Unlock()

	supplier, exists := suppliers[req.Id]
	if !exists {
		return nil, errors.New("supplier not found")
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.Contact != "" {
		supplier.Contact = req.Contact
	}

	log.Println("Supplier updated:", supplier)
	return supplier, nil
}

// Delete a supplier
func (s *SupplierServiceProto) P_DeleteSupplier(ctx context.Context, req *DeleteSupplierRequest) (*DeleteSupplierResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	_, exists := suppliers[req.Id]
	if !exists {
		return nil, errors.New("supplier not found")
	}

	delete(suppliers, req.Id)

	log.Println("Supplier deleted:", req.Id)
	return &DeleteSupplierResponse{Message: "Supplier deleted successfully"}, nil
}
