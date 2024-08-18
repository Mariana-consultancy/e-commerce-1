package repository

import "e-commerce/internal/models"

// Function to find all the items in the database
func (p *Postgres) FindAllItems() ([]models.Item, error) {
	// Init a variable called items using an array of the models.Item data type (Models Folder)
	// An array because there's multiple items we'll need
	var items []models.Item

	// init a variable that error tests whether the variable exists in the database 
	if err := p.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	// return the variable if found without errors
	return items, nil
}

// Function to get items by their unique ID
func (p *Postgres) GetItemsByID(itemID uint) (*models.Item, error) {
	// init a variable that contains a slice of the models.item - it's a pointer because we need to use the data
	item := &models.Item{}

	// init an error test that checks if the item exists by ID
	if err := p.DB.Where("ID = ?", itemID).First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Function to create an item in the database
func (p *Postgres) CreateItem(item *models.Item) error {
	// init an error test on creating an item in the database
	if err := p.DB.Create(item).Error; err != nil {
		return err
	}
	return nil
}

// Update a user in the database
func (p *Postgres) UpdateItem(item *models.Item) error {
	// init an error test on updating an item in the database
	if err := p.DB.Save(item).Error; err != nil {
		return err
	}
	return nil
}
