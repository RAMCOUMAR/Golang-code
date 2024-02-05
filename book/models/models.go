// models/models.go
package models

// User represents the user entity with additional details
type User struct {
	UserID       int    `json:"userid"`
	UserName     string `json:"username"`
	AdminID      int    `json:"admin_id"`
	Role         string `json:"role"`
	RoleID       int    `json:"role_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	SessionToken string `json:"session_token"`
}

// Product represents the product entity
type Product struct {
	ProductID int    `json:"productid"`
	Name      string `json:"name"`
	Price     string `json:"price"`
}

// ProductUser represents a combined entity of product and user details
type ProductUser struct {
	ProductID    int    `json:"productid"`
	ProductName  string `json:"name"`
	ProductPrice string `json:"price"`
	UserID       int    `json:"userid"`
	UserName     string `json:"username"`
}
