package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	Order_ID
	Order_Cart
	Ordered_At
	Price
	Discount
	Payment_Method
}
