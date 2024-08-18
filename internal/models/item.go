package models

import (
	"gorm.io/gorm"
)

// Create a struct data type for the object

type Item struct {
	gorm.Model
	Id       uint   `json:"id"`
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Total    int    `json:"total"`
	Discount int    `json:"discount"`
}
