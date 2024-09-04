package repository

import (
	"e-commerce/internal/models"
)

func (p *Postgres) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Postgres) FindAllUsers() ([]models.User, error) {
	var users []models.User

	if err := p.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (p *Postgres) GetUserByID(userID uint) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("ID = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Postgres) CreateUser(user *models.User) error {
	if err := p.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateUser(user *models.User) error {
	if err := p.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetProductByID(productID uint) (*models.Product, error) {
	product := &models.Product{}

	if err := p.DB.Where("ID = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Postgres) GetAllProducts() ([]models.Product, error) {
	var products []models.Product

	if err := p.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *Postgres) AddToCart(cart *models.Cart) error {
	if err := p.DB.Save(cart).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetCartItemByProductID(productID uint) (*models.Cart, error) {
	cart := &models.Cart{}

	if err := p.DB.Where("ID = ?", productID).First(&cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (p *Postgres) DeleteProductFromCart(cart *models.Cart) error {
	if err := p.DB.Delete(cart).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetCartByUserID(userID uint) ([]models.Cart, error) {
	var cartinventory []models.Cart

	if err := p.DB.Where("ID = ?", userID).First(&cartinventory).Error; err != nil {
		return nil, err

	}
	return cartinventory, nil

}

func (p *Postgres) GetOrderItemsByOrderID(orderID uint) ([]*models.OrderItem, error) {
	var orderDetails []*models.OrderItem

	if err := p.DB.Preload("Product").Where("order_id = ?", orderID).Find(&orderDetails).Error; err != nil {
		return nil, err
	}
	return orderDetails, nil
}
