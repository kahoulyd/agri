package models

type Supplier struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}
