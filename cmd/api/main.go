package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/geoffreyoliaro/mininexus/internal/graph"
	"github.com/geoffreyoliaro/mininexus/internal/handlers"
	"github.com/geoffreyoliaro/mininexus/internal/middleware"
)

func main() {
	// Load .env in local dev; Cloud Run injects env vars directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading env from environment")
	}

	// Connect to FalkorDB (supports both local and cloud deployments)
	// Cloud: uses FALKORDB_USER and FALKORDB_PASSWORD for authentication
	// Local: uses FALKORDB_PASSWORD only (can be empty)
	db, err := graph.NewClient(
		os.Getenv("FALKORDB_HOST"),
		os.Getenv("FALKORDB_PORT"),
		os.Getenv("FALKORDB_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("FalkorDB connection failed: %v", err)
	}
	defer db.Close()

	// Seed demo graph on startup (idempotent)
	if err := graph.SeedDemoData(db); err != nil {
		log.Printf("Seed warning: %v", err)
	}

	r := gin.Default()

	// Health check — Cloud Run uses this
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public: Xano validates credentials, we issue our own JWT
	r.POST("/auth/login", handlers.Login)

	// Protected routes — JWT issued by this service after Xano validation
	api := r.Group("/api/v1", middleware.JWTAuth())
	{
		// People
		api.GET("/people", handlers.ListPeople(db))
		api.POST("/people", handlers.CreatePerson(db))
		api.GET("/people/:id", handlers.GetPerson(db))

		// Relationships
		api.POST("/relationships", handlers.CreateRelationship(db))
		api.GET("/relationships/:id/network", handlers.GetNetwork(db))

		// Intelligence queries — the "money" endpoints for the demo
		api.GET("/intelligence/mutual-connections", handlers.MutualConnections(db))
		api.GET("/intelligence/shortest-path", handlers.ShortestPath(db))
		api.GET("/intelligence/relationship-strength", handlers.RelationshipStrength(db))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("MiniNexus API listening on :%s", port)
	r.Run(":" + port)
}
