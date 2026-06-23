package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/geoffreyoliaro/mininexus/internal/graph"
	"github.com/geoffreyoliaro/mininexus/internal/models"
)

// ListPeople returns all Person nodes from the graph.
func ListPeople(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := db.ListPeople(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		people := []models.Person{}
		for result.Next() {
			r := result.Record()
			people = append(people, models.Person{
				ID:      str(r.GetByIndex(0)),
				Name:    str(r.GetByIndex(1)),
				Email:   str(r.GetByIndex(2)),
				Company: str(r.GetByIndex(3)),
				Role:    str(r.GetByIndex(4)),
			})
		}
		c.JSON(http.StatusOK, gin.H{"data": people, "count": len(people)})
	}
}

// CreatePerson creates a new Person node.
func CreatePerson(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreatePersonRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := uuid.New().String()
		if err := db.CreatePerson(c.Request.Context(), id, req.Name, req.Email, req.Company, req.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": models.Person{
				ID:      id,
				Name:    req.Name,
				Email:   req.Email,
				Company: req.Company,
				Role:    req.Role,
			},
		})
	}
}

// GetPerson returns a person and their 1-hop relationship network.
func GetPerson(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result, err := db.GetPersonNetwork(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		edges := []models.NetworkEdge{}
		for result.Next() {
			r := result.Record()
			edges = append(edges, models.NetworkEdge{
				Source:        str(r.GetByIndex(0)),
				Target:        str(r.GetByIndex(1)),
				TargetCompany: str(r.GetByIndex(2)),
				RelType:       str(r.GetByIndex(3)),
				Strength:      intVal(r.GetByIndex(4)),
			})
		}
		c.JSON(http.StatusOK, gin.H{"person_id": id, "network": edges})
	}
}

// CreateRelationship adds a KNOWS edge between two people.
func CreateRelationship(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateRelationshipRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.CreateRelationship(c.Request.Context(), req.FromID, req.ToID, req.Type, req.Strength); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "relationship created", "data": req})
	}
}
