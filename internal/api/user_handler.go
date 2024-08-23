package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.Password = hashedPassword

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}

// Login User
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequestUser
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Email and Password must not be empty", 400, nil, nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	// Verify the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating access token", 500, err.Error(), nil)
		return
	}

	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating refresh token", 500, err.Error(), nil)
		return
	}

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// View Product listing -
func (u *HTTPHandler) GetAllProducts(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	products, err := u.Repository.GetAllProducts()
	if err != nil {
		util.Response(c, "Error getting products", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, products, nil)
}

// View Product by ID
func (u *HTTPHandler) GetProductByID(c *gin.Context) {
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	productID := c.Param("id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	product, err := u.Repository.GetProductByID(uint(id))
	if err != nil {
		util.Response(c, "Error getting product", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Success", 200, product, nil)
}

// Add to Cart
func (u *HTTPHandler) AddToCart(c *gin.Context) {
	// Authorisation for Cart -
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}
	// Create a variable for to store the item's data using the Item model
	var cart *models.Cart
	// initialise err variable, binds the request data to the item struct
	if err := c.ShouldBind(&cart); err != nil {
		// Return an error message if there's an error
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}
	// Get product by ID
	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 401, err.Error(), nil)
		return
	}
	// Check cart quantity
	if cart.Quantity > product.Quantity {
		util.Response(c, "Not enough products ", 400, nil, nil)
		return
	}
	cart.UserID = user.ID
	// Init an err variable connected to the repository interface the newly created item
	err = u.Repository.AddToCart(cart)
	// Implement an error test to see if the item was successfully created
	if err != nil {
		util.Response(c, "Cart not created", 500, err.Error(), nil)
		return
	}
	// return a successfull message if the item was created
	util.Response(c, "Cart created", 200, nil, nil)
}

// Edit Cart - think like a thief (break the system)
func (u *HTTPHandler) EditCart(c *gin.Context) {
	// Authorisation for Cart -
	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}
	// Request Body
	var cart *models.Cart
	err = c.ShouldBind(&cart)
	if err != nil {
		// Return an error message if there's an error
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Make sure the cart exists
	shoppingCart, err := u.Repository.GetCartItemByProductID(cart.ProductID)
	if err != nil {
		util.Response(c, "Cart not found", 500, err.Error(), nil)
		return
	}

	// Ensure product exist
	product, err := u.Repository.GetProductByID(cart.ProductID)
	if err != nil {
		util.Response(c, "Product not found", 500, err.Error(), nil)
		return
	}

	// Check cart quantity
	if cart.Quantity > product.Quantity {
		util.Response(c, "Not enough products ", 400, nil, nil)
		return
	}

	cart.UserID = user.ID
	cart.ID = shoppingCart.ID

	// Add to cart db method
	err = u.Repository.AddToCart(cart)
	if err != nil {
		util.Response(c, "Error editing product quantity", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Product quantity edited ", 200, nil, nil)
}

// RemoveFromCart removes a product from the user's cart
func (u *HTTPHandler) RemoveFromCart(c *gin.Context) {
	// Authorize the user
	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	var cart *models.Cart
	// initialise err variable, binds the request data to the item struct
	if err := c.ShouldBind(&cart); err != nil {
		// Return an error message if there's an error
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	// Get the product ID from the URL parameters
	productIDParam := c.Param("id")
	if productIDParam == "" {
		util.Response(c, "Product ID is required", 400, nil, nil)
		return
	}
	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		util.Response(c, "invalid product ID", 400, err.Error(), nil)
		return
	}

	shoppingCart, err := u.Repository.GetCartItemByProductID(uint(productID))
	if err != nil {
		util.Response(c, "Product not found", 404, err.Error(), nil)
		return
	}

	// Remove the product from the cart
	err = u.Repository.DeleteProductFromCart(shoppingCart)
	if err != nil {
		util.Response(c, "Error removing product from cart", 500, err.Error(), nil)
		return
	}

	// Return a success response
	util.Response(c, "Product removed from cart", 200, nil, nil)
}

// view cart 
func (u *HTTPHandler) ViewCart(c *gin.Context) {
	user, err := u.GetUserFromContext(c)
	if err !=nil{



		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

		cartinventory, err := u.Repository.GetCartByUserID(user.ID)
		if err != nil {
			util.Response(c, "Error Dsiplaying Cart", 500, err.Error(), nil)
			return
		}
  util.Response(c, "Product added to cart", 200, cartinventory, nil)
}

