// handlers.go
package handlers

import (
	"book/jwt"
	"book/models"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

// InitializeDB initializes the database connection
func InitializeDB(database *sql.DB) {
	db = database
}

// getProductsHandler fetches all products from the database
func GetProductsHandler(c *gin.Context) {
	query := "SELECT productid, name, price FROM product"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []models.Product

	for rows.Next() {
		var result models.Product
		err := rows.Scan(&result.ProductID, &result.Name, &result.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
}

// insertUserHandler handles the insertion of user data
func InsertUserHandler(c *gin.Context) {
	var data models.User
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertStatement := "INSERT INTO user (userid, username, admin_id, role, role_id, first_name, last_name, session_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(insertStatement, data.UserID, data.UserName, data.AdminID, data.Role, data.RoleID, data.FirstName, data.LastName, data.SessionToken)
	if err != nil {
		log.Println("Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Data Received successfully"})
}

// getUsersHandler fetches all users from the database
func GetUsersHandler(c *gin.Context) {
	query := "SELECT userid, username, admin_id, role, role_id, first_name, last_name, session_token FROM user"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []models.User

	for rows.Next() {
		var result models.User
		err := rows.Scan(&result.UserID, &result.UserName, &result.AdminID, &result.Role, &result.RoleID, &result.FirstName, &result.LastName, &result.SessionToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
}

// loginHandler handles user login and generates JWT token
func LoginHandler(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authenticateUser(user) {
		token, err := jwt.GenerateToken(user.UserID, user.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		userDetails := getUserDetails(user.UserID)

		response := gin.H{
			"code":    http.StatusOK,
			"message": "User logged in successfully",
			"result": gin.H{
				"admin_id":      userDetails.AdminID,
				"first_name":    userDetails.FirstName,
				"last_name":     userDetails.LastName,
				"role":          userDetails.Role,
				"role_id":       userDetails.RoleID,
				"session_token": userDetails.SessionToken,
				"username":      userDetails.UserName,
				"token":         token,
			},
		}

		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// getUserDetails retrieves additional user details from the database
func getUserDetails(userID int) *models.User {
	query := `
        SELECT admin_id, first_name, last_name, role, role_id, session_token, username
        FROM user
        WHERE userid = ?
    `

	var user models.User
	err := db.QueryRow(query, userID).Scan(
		&user.AdminID,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.RoleID,
		&user.SessionToken,
		&user.UserName,
	)
	if err != nil {
		return nil
	}

	user.UserID = userID

	return &user
}

// authenticateUser checks user credentials against a database
func authenticateUser(user models.User) bool {
	query := "SELECT COUNT(*) FROM user WHERE userid = ? AND username = ?"
	var count int
	err := db.QueryRow(query, user.UserID, user.UserName).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// getUserClaims retrieves user claims from the JWT token
func GetUserClaims(c *gin.Context) (int, string) {
	claims, _ := c.Get("claims")
	userID := claims.(*jwt.Claims).UserID
	username := claims.(*jwt.Claims).Username
	return userID, username
}

// insertProductHandler handles the insertion of product data
func InsertProductHandler(c *gin.Context) {
	var data models.Product
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertStatement := "INSERT INTO product (productid, name, price) VALUES (?, ?, ?)"
	_, err := db.Exec(insertStatement, data.ProductID, data.Name, data.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Data Received successfully"})
}

// updateProductHandler handles the update of product data
func UpdateProductHandler(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productID"})
		return
	}

	var data models.Product
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateStatement := "UPDATE product SET name = ?, price = ? WHERE productid = ?"
	_, err = db.Exec(updateStatement, data.Name, data.Price, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Data Updated successfully"})
}

// deleteProductHandler handles the deletion of product data
func DeleteProductHandler(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productID"})
		return
	}

	deleteStatement := "DELETE FROM product WHERE productid = ?"
	_, err = db.Exec(deleteStatement, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Data Deleted successfully"})
}

// fetchDataHandler fetches data using INNER JOIN from the database
func FetchDataHandler(c *gin.Context) {
	query := `
        SELECT p.productid, p.name, p.price, u.userid, u.username
        FROM product p
        INNER JOIN user u ON p.productid = u.userid
    `

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []models.ProductUser

	for rows.Next() {
		var result models.ProductUser
		err := rows.Scan(&result.ProductID, &result.ProductName, &result.ProductPrice, &result.UserID, &result.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
}
