package main

import (
	"e-commerce/cmd/server"
	"e-commerce/internal/repository"
)

func main() {
	//Gets the environment variables
	env := server.InitDBParams()

	//Initializes the database
	db, err := repository.Initialize(env.DbUrl)
	if err != nil {
		return
	}

	//Runs the app
	server.Run(db, env.Port)
}

// application

// database

// Create   Read    Update   Delete

//HTTP Methods
//C = Create : storing information/data in the backend  e.g sign up, create a bank account etc [POST] [Request Body]
//R = Read	  : getting information/data from the backend e.g viewing profiles [GET] Params [query parameter and path parameter]
//U = Update  :  updating an information in the backend  e.g changing of information(age, name etc) [PUT, PATCH] [Request Body]
//D = Delete  : Deleting an information from the backend  (soft delete and hard delete) e.g remove you account from a platform . [DELETE]

//Authenticated Routes : It gets information that is related to the user

//Endpoints to Create for an E-commerce Application
//1. Create products [POST] //  Seller Authenticated
//2. Get list of products [GET] ?filter=fashion   Not authenticated
//3. Create an account [POST]  // buyer and seller  Not Authenticated
//4. Log in [POST]   // buyer login & seller login   Not authenticated
//5. Add to cart [PUT / PATCH]   Buyer Authenticated
//6. Remove From Cart [PUT/PATCH]  Buyer Authenticated
//7. Add to wishlist  [PUT / PATCH]   Buyer Authenticated
//8. Remove from wishlist  [PUT / PATCH]   Buyer Authenticated
// 9.