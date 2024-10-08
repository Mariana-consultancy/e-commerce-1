package repository

import "e-commerce/internal/models"

func (p *Postgres) FindSellerByEmail(email string) (*models.Seller, error) {
	seller := &models.Seller{}

	if err := p.DB.Where("email = ?", email).First(&seller).Error; err != nil {
		return nil, err
	}
	return seller, nil
}

// Create a seller in the database
func (p *Postgres) CreateSeller(seller *models.Seller) error {
	if err := p.DB.Create(seller).Error; err != nil {
		return err
	}
	return nil
}

// Update a user in the database
func (p *Postgres) UpdateSeller(seller *models.Seller) error {
	if err := p.DB.Save(seller).Error; err != nil {
		return err
	}
	return nil
}

// Create product
func (p *Postgres) CreateProduct(product *models.Product) error {
	if err := p.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetOrderByID(orderID uint) (*models.Order, error) {
	order := &models.Order{}

	if err := p.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (p *Postgres) UpdateOrder(order *models.Order) error {
	if err := p.DB.Save(order).Error; err != nil {
		return err
	}
	return nil
}
