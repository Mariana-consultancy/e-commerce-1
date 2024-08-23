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

func (p *Postgres) GetCartByUserID(userID uint) (*models.Cart, error) {
	cart := &models.Cart{}
	if err := p.DB.Where("ID = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (p *Postgres) RemoveFromCart(userID uint, productID uint) error {
	// First, find the cart associated with the user ID
	cart := &models.Cart{}
	if err := p.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}

	// Find the cart item to remove based on the cart ID and product ID
	cartItem := &models.Cart{}
	if err := p.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).Delete(&cartItem).Error; err != nil {
		return err
	}

	// Delete the cart item
	if err := p.DB.Delete(cartItem).Error; err != nil {
		return err
	}

	return nil
}
