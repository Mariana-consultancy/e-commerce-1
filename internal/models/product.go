package models

import ("gorm.io/gorm")

type Product struct {
	gorm.Model
	SellerID		uint		`json:""`
	Name			string		`json:""`
	Price			float64		`json:""`
	ImageUrl		string		`json:""`
	Quantity		int			`json:""`
	Description		string		`json:""`
	Status			bool		`json:""`
}