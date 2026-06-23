package graph

import (
	"context"
	"fmt"
	"strconv"

	falkordb "github.com/FalkorDB/falkordb-go"
)

// Client wraps a FalkorDB graph connection.
type Client struct {
	db    *falkordb.FalkorDB
	graph *falkordb.Graph
}

// NewClient connects to FalkorDB (Redis-compatible protocol).
func NewClient(host, port, password string) (*Client, error) {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid FALKORDB_PORT: %w", err)
	}

	// Support both local (no auth) and cloud (with auth) deployments
	db, err := falkordb.FalkorDBNew(&falkordb.ConnectionOption{
		Host:     host,
		Port:     p,
		Password: password,
		// Username is optional; used for FalkorDB Cloud authentication
	})
	if err != nil {
		return nil, err
	}

	g := db.SelectGraph("nexus")
	return &Client{db: &db, graph: &g}, nil
}

// Close tears down the connection.
func (c *Client) Close() {
	c.db.Close()
}

// Query executes a raw Cypher query and returns result set.
func (c *Client) Query(ctx context.Context, cypher string, params map[string]interface{}) (*falkordb.QueryResult, error) {
	result, err := c.graph.QueryRO(cypher, params)
	if err != nil {
		// Retry as write query if read-only fails
		result, err = c.graph.Query(cypher, params)
	}
	return result, err
}

// WriteQuery executes a write Cypher query.
func (c *Client) WriteQuery(ctx context.Context, cypher string, params map[string]interface{}) (*falkordb.QueryResult, error) {
	return c.graph.Query(cypher, params)
}

// ---- Graph schema operations ----

// CreatePerson creates a Person node.
// Cypher: MERGE prevents duplicates on re-run.
func (c *Client) CreatePerson(ctx context.Context, id, name, email, company, role string) error {
	cypher := `
		MERGE (p:Person {id: $id})
		SET p.name = $name,
		    p.email = $email,
		    p.company = $company,
		    p.role = $role
		RETURN p`
	_, err := c.WriteQuery(ctx, cypher, map[string]interface{}{
		"id":      id,
		"name":    name,
		"email":   email,
		"company": company,
		"role":    role,
	})
	return err
}

// CreateRelationship creates a directed KNOWS edge with a strength score.
func (c *Client) CreateRelationship(ctx context.Context, fromID, toID, relType string, strength int) error {
	cypher := `
		MATCH (a:Person {id: $from}), (b:Person {id: $to})
		MERGE (a)-[r:KNOWS {type: $type}]->(b)
		SET r.strength = $strength,
		    r.updated_at = timestamp()
		RETURN r`
	_, err := c.WriteQuery(ctx, cypher, map[string]interface{}{
		"from":     fromID,
		"to":       toID,
		"type":     relType,
		"strength": strength,
	})
	return err
}

// GetPersonNetwork returns a person and their 1-hop connections.
func (c *Client) GetPersonNetwork(ctx context.Context, personID string) (*falkordb.QueryResult, error) {
	cypher := `
		MATCH (p:Person {id: $id})-[r:KNOWS]-(connected:Person)
		RETURN p.name AS source,
		       connected.name AS target,
		       connected.company AS target_company,
		       r.type AS rel_type,
		       r.strength AS strength
		ORDER BY r.strength DESC`
	return c.Query(ctx, cypher, map[string]interface{}{"id": personID})
}

// MutualConnections finds people both A and B know.
func (c *Client) MutualConnections(ctx context.Context, idA, idB string) (*falkordb.QueryResult, error) {
	cypher := `
		MATCH (a:Person {id: $idA})-[:KNOWS]-(mutual:Person)-[:KNOWS]-(b:Person {id: $idB})
		WHERE a <> b
		RETURN mutual.name AS name,
		       mutual.company AS company,
		       mutual.role AS role`
	return c.Query(ctx, cypher, map[string]interface{}{"idA": idA, "idB": idB})
}

// ShortestPath finds the shortest connection chain between two people.
// This is FalkorDB's graph power on full display.
func (c *Client) ShortestPath(ctx context.Context, fromID, toID string) (*falkordb.QueryResult, error) {
	cypher := `
		MATCH (a:Person {id: $from}), (b:Person {id: $to}),
		      p = shortestPath((a)-[:KNOWS*..6]-(b))
		RETURN [node IN nodes(p) | node.name] AS path,
		       length(p) AS hops`
	return c.Query(ctx, cypher, map[string]interface{}{"from": fromID, "to": toID})
}

// RelationshipStrength computes a weighted score between two people
// based on shared connections and direct edge strength.
func (c *Client) RelationshipStrength(ctx context.Context, idA, idB string) (*falkordb.QueryResult, error) {
	cypher := `
		MATCH (a:Person {id: $idA}), (b:Person {id: $idB})
		OPTIONAL MATCH (a)-[direct:KNOWS]-(b)
		OPTIONAL MATCH (a)-[:KNOWS]-(mutual:Person)-[:KNOWS]-(b)
		WITH a, b,
		     COALESCE(direct.strength, 0) AS direct_strength,
		     COUNT(DISTINCT mutual) AS mutual_count
		RETURN a.name AS person_a,
		       b.name AS person_b,
		       direct_strength,
		       mutual_count,
		       (direct_strength * 0.6 + mutual_count * 10) AS composite_score`
	return c.Query(ctx, cypher, map[string]interface{}{"idA": idA, "idB": idB})
}

// ListPeople returns all Person nodes.
func (c *Client) ListPeople(ctx context.Context) (*falkordb.QueryResult, error) {
	cypher := `
		MATCH (p:Person)
		RETURN p.id AS id, p.name AS name, p.email AS email,
		       p.company AS company, p.role AS role
		ORDER BY p.name`
	return c.Query(ctx, cypher, nil)
}
