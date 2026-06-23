package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/geoffreyoliaro/mininexus/internal/graph"
	"github.com/geoffreyoliaro/mininexus/internal/models"
)

// MutualConnections returns people that both person A and B know.
// GET /api/v1/intelligence/mutual-connections?a=p1&b=p4
func MutualConnections(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idA := c.Query("a")
		idB := c.Query("b")
		if idA == "" || idB == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query params 'a' and 'b' are required"})
			return
		}

		result, err := db.MutualConnections(c.Request.Context(), idA, idB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mutuals := []models.MutualConnection{}
		for result.Next() {
			r := result.Record()
			mutuals = append(mutuals, models.MutualConnection{
				Name:    str(r.GetByIndex(0)),
				Company: str(r.GetByIndex(1)),
				Role:    str(r.GetByIndex(2)),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"person_a":            idA,
			"person_b":            idB,
			"mutual_connections":  mutuals,
			"count":               len(mutuals),
			"cypher_hint": "MATCH (a)-[:KNOWS]-(m:Person)-[:KNOWS]-(b) WHERE a<>b RETURN m",
		})
	}
}

// ShortestPath finds the shortest chain of KNOWS edges between two people.
// GET /api/v1/intelligence/shortest-path?from=p1&to=p8
func ShortestPath(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromID := c.Query("from")
		toID := c.Query("to")
		if fromID == "" || toID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query params 'from' and 'to' are required"})
			return
		}

		result, err := db.ShortestPath(c.Request.Context(), fromID, toID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !result.Next() {
			c.JSON(http.StatusNotFound, gin.H{"error": "no path found between these people"})
			return
		}

		r := result.Record()
		pathRaw := r.GetByIndex(0)
		hops := intVal(r.GetByIndex(1))

		// Convert []interface{} to []string
		path := []string{}
		if nodes, ok := pathRaw.([]interface{}); ok {
			for _, n := range nodes {
				path = append(path, fmt.Sprintf("%v", n))
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"from": fromID,
			"to":   toID,
			"path": models.PathResult{Path: path, Hops: hops},
			"cypher_hint": "shortestPath((a)-[:KNOWS*..6]-(b))",
		})
	}
}

// RelationshipStrength computes a composite score between two people.
// GET /api/v1/intelligence/relationship-strength?a=p1&b=p3
func RelationshipStrength(db *graph.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idA := c.Query("a")
		idB := c.Query("b")
		if idA == "" || idB == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query params 'a' and 'b' are required"})
			return
		}

		result, err := db.RelationshipStrength(c.Request.Context(), idA, idB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !result.Next() {
			c.JSON(http.StatusNotFound, gin.H{"error": "could not compute strength"})
			return
		}

		r := result.Record()
		score := models.StrengthResult{
			PersonA:        str(r.GetByIndex(0)),
			PersonB:        str(r.GetByIndex(1)),
			DirectStrength: intVal(r.GetByIndex(2)),
			MutualCount:    intVal(r.GetByIndex(3)),
			CompositeScore: floatVal(r.GetByIndex(4)),
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    score,
			"formula": "composite = (direct_strength × 0.6) + (mutual_count × 10)",
		})
	}
}

// ---- helpers ----

func str(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func intVal(v interface{}) int {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case int64:
		return int(t)
	case float64:
		return int(t)
	}
	return 0
}

func floatVal(v interface{}) float64 {
	if v == nil {
		return 0
	}
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}
