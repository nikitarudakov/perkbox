package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nikitarudakov/perkbox/internal/handlers"
	"github.com/nikitarudakov/perkbox/internal/repo"
	"log"
)

func main() {
	cfg, err := repo.LoadConfig()
	if err != nil {
		log.Fatalln("error loading repo config: ", err)
	}

	repos, err := repo.NewRepository(cfg)
	if err != nil {
		log.Fatalln("error creating repository: ", err)
	}

	userHandler := handlers.NewUserHandler(repos)

	r := gin.Default()

	// Middleware to simulate admin-only access
	adminRoutes := r.Group("/api", RequireRole("admin"))
	adminRoutes.POST("/users", userHandler.CreateUser)
	adminRoutes.DELETE("/users/:user_id", userHandler.DeleteUser)
	adminRoutes.GET("/businesses/:business_id/users", userHandler.ListUsers)

	// Routes available to all users in the same business
	r.PUT("/api/users/:user_id", userHandler.UpdateUser)
	r.GET("/api/users/:user_id", userHandler.GetUser)

	log.Println("Server running on :8080")
	r.Run(":8080")
}

// RequireRole checks if the user has the required role via header.
func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		if role != required {
			c.AbortWithStatusJSON(403, gin.H{"error": "admin role required"})
			return
		}
		c.Next()
	}
}
