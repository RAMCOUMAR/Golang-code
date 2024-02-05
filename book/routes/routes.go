// routes/routes.go
package routes

import (
	"book/database"
	"book/handlers"
	"book/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all routes
func SetupRoutes(r *gin.Engine, db *database.Database) {
	// Public routes (no authentication required)
	r.POST("/api/login", handlers.LoginHandler)

	// Protected route (require JWT authentication)
	protected := r.Group("/api")
	protected.Use(middlewares.JWTMiddleware())
	{
		protected.POST("/insert-product", handlers.InsertProductHandler)
	}

	// Public routes continued
	r.GET("/api/get-products", handlers.GetProductsHandler)
	r.POST("/api/insert-user", handlers.InsertUserHandler)
	r.GET("/api/get-users", handlers.GetUsersHandler)
	r.PUT("/api/update-product/:productID", handlers.UpdateProductHandler)
	r.DELETE("/api/delete-product/:productID", handlers.DeleteProductHandler)
	r.GET("/api/fetch-data", handlers.FetchDataHandler)
}
