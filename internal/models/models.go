package models

// Person represents a node in the relationship graph.
type Person struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Role    string `json:"role"`
}

// Relationship represents a KNOWS edge between two people.
type Relationship struct {
	FromID   string `json:"from_id"`
	ToID     string `json:"to_id"`
	Type     string `json:"type"`    // investor, advisor, colleague, etc.
	Strength int    `json:"strength"` // 0-100
}

// NetworkEdge is a single edge returned from a network query.
type NetworkEdge struct {
	Source        string `json:"source"`
	Target        string `json:"target"`
	TargetCompany string `json:"target_company"`
	RelType       string `json:"rel_type"`
	Strength      int    `json:"strength"`
}

// MutualConnection represents a shared contact.
type MutualConnection struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Role    string `json:"role"`
}

// PathResult is the result of a shortest path query.
type PathResult struct {
	Path []string `json:"path"`
	Hops int      `json:"hops"`
}

// StrengthResult is the result of a relationship strength computation.
type StrengthResult struct {
	PersonA        string  `json:"person_a"`
	PersonB        string  `json:"person_b"`
	DirectStrength int     `json:"direct_strength"`
	MutualCount    int     `json:"mutual_count"`
	CompositeScore float64 `json:"composite_score"`
}

// CreatePersonRequest is the inbound payload for creating a person.
type CreatePersonRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Company string `json:"company" binding:"required"`
	Role    string `json:"role" binding:"required"`
}

// CreateRelationshipRequest is the inbound payload for creating a relationship.
type CreateRelationshipRequest struct {
	FromID   string `json:"from_id" binding:"required"`
	ToID     string `json:"to_id" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Strength int    `json:"strength" binding:"min=1,max=100"`
}
