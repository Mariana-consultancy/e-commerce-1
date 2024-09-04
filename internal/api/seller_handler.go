package api

import (
	"e-commerce/internal/middleware"
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Create Seller
func (u *HTTPHandler) CreateSeller(c *gin.Context) {
	var seller *models.Seller
	if err := c.ShouldBind(&seller); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	isEmailExist, _ := u.Repository.FindSellerByEmail(seller.Email)
	if isEmailExist != nil {
		util.Response(c, "Email already exist", 400, nil, nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(seller.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	seller.Password = hashedPassword

	err = u.Repository.CreateSeller(seller)
	if err != nil {
		util.Response(c, "Seller not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Seller created", 200, nil, nil)

}

// Login Seller
func (u *HTTPHandler) LoginSeller(c *gin.Context) {
	var loginRequest *models.LoginRequestSeller
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		util.Response(c, "Email and Password must not be empty", 400, nil, nil)
		return
	}

	seller, err := u.Repository.FindSellerByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	// Verify the password
	if err = bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(loginRequest.Password)); err != nil {
		util.Response(c, "invalid email or password", 400, "invalid email or password", nil)
		return
	}

	accessClaims, refreshClaims := middleware.GenerateClaims(seller.Email)

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
		"seller":        seller,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// Create Product
func (u *HTTPHandler) CreateProduct(c *gin.Context) {
	seller, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "invalid token", 401, err.Error(), nil)
		return
	}

	var product *models.Product
	if err := c.ShouldBind(&product); err != nil {
		util.Response(c, "invalid request", 401, err.Error(), nil)
		return
	}

	product.SellerID = seller.ID
	product.Status = false

	err = u.Repository.CreateProduct(product)
	if err != nil {
		util.Response(c, "product not created", 500, err.Error(), nil)
		return
	}

	util.Response(c, "product created successfully", 201, nil, nil)

}

func (u *HTTPHandler) AcceptOrder(c *gin.Context) {
	_, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		util.Response(c, "Order ID not provided", 400, nil, nil)
		return
	}

	//convert id to uint
	orderIDUint, err := util.ConvertStringToUint(orderID)
	if err != nil {
		util.Response(c, "Invalid order ID", 400, err.Error(), nil)
		return
	}

	// Get the order from the database
	order, err := u.Repository.GetOrderByID(orderIDUint)
	if err != nil {
		util.Response(c, "Order not found", 404, err.Error(), nil)
		return
	}

	//check if order is already accepted
	if order.Status == "ACCEPTED" {
		util.Response(c, "Order already accepted", 400, nil, nil)
		return
	}

	// Update the order status to accepted
	order.Status = "ACCEPTED"
	if err := u.Repository.UpdateOrder(order); err != nil {
		util.Response(c, "Error updating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order accepted", 200, nil, nil)
}

// Decline the order
func (u *HTTPHandler) DeclineOrder(c *gin.Context) {
	_, err := u.GetSellerFromContext(c)
	if err != nil {
		util.Response(c, "Invalid token", 401, err.Error(), nil)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		util.Response(c, "Order ID not provided", 400, nil, nil)
		return
	}

	//convert id to uint
	orderIDUint, err := util.ConvertStringToUint(orderID)
	if err != nil {
		util.Response(c, "Invalid order ID", 400, err.Error(), nil)
		return
	}

	// Get the order from the database
	order, err := u.Repository.GetOrderByID(orderIDUint)
	if err != nil {
		util.Response(c, "Order not found", 404, err.Error(), nil)
		return
	}

	//check if order is already declined
	if order.Status == "DECLINED" {
		util.Response(c, "Order already declined", 400, nil, nil)
		return
	}

	// Update the order status to declined
	order.Status = "DECLINED"
	if err := u.Repository.UpdateOrder(order); err != nil {
		util.Response(c, "Error updating order", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Order declined", 200, nil, nil)
}
