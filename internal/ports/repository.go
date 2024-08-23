package ports

import "e-commerce/internal/models"

type Repository interface {
	// Find Functions
	FindUserByEmail(email string) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	FindSellerByEmail(email string) (*models.Seller, error)
	// Create Functions
	CreateUser(user *models.User) error
	CreateSeller(Seller *models.Seller) error
	CreateProduct(product *models.Product) error
	// Get Functions
	GetCartItemByProductID(productID uint) (*models.Cart, error)
	GetUserByID(userID uint) (*models.User, error)
	GetProductByID(productID uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	// Update Functions
	UpdateUser(user *models.User) error
	UpdateSeller(user *models.Seller) error
	// Add Functions
	AddToCart(cart *models.Cart) error
	// Tokens
	BlacklistToken(token *models.BlacklistTokens) error
	TokenInBlacklist(token *string) bool
	// Remove From Cart
	DeleteProductFromCart(cart *models.Cart) error
}
