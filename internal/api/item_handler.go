package api

import (
	"e-commerce/internal/models"
	"e-commerce/internal/util"
	"github.com/gin-gonic/gin"
)

// Create Item function
func (u *HTTPHandler) CreateItem(c *gin.Context) {
	// Create a variable for to store the item's data using the Item model 
	var item *models.Item

	// initialise err variable, binds the request data to the item struct
	if err := c.ShouldBind(&item); err != nil {
		// Return an error message if there's an error
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}
	// Init an err variable connected to the repository interface the newly created item 
	err := u.Repository.CreateItem(item)

	// Implement an error test to see if the item was successfully created
	if err != nil {
		util.Response(c, "Item not created", 500, err.Error(), nil)
		return
	}
	// return a successfull message if the item was created
	util.Response(c, "Item created", 200, nil, nil)
}
