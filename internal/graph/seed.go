package graph

import (
	"context"
	"log"
)

// SeedDemoData populates the graph with a realistic relationship network.
// Uses MERGE so it's safe to call on every startup.
func SeedDemoData(c *Client) error {
	ctx := context.Background()
	log.Println("Seeding demo graph data...")

	// People — a mix of founders, investors, operators
	people := []struct{ id, name, email, company, role string }{
		{"p1", "Amara Osei", "amara@horizonvc.com", "Horizon VC", "Partner"},
		{"p2", "Lena Fischer", "lena@buildco.io", "BuildCo", "CEO"},
		{"p3", "James Mwangi", "james@techbridge.ke", "TechBridge", "CTO"},
		{"p4", "Sofia Reyes", "sofia@openloop.ai", "OpenLoop AI", "Founder"},
		{"p5", "Tom Nakamura", "tom@horizonvc.com", "Horizon VC", "Analyst"},
		{"p6", "Priya Sharma", "priya@buildco.io", "BuildCo", "VP Engineering"},
		{"p7", "David Oliaro", "david@nexustech.io", "NexusTech", "Head of Growth"},
		{"p8", "Chioma Eze", "chioma@openloop.ai", "OpenLoop AI", "COO"},
	}
	for _, p := range people {
		if err := c.CreatePerson(ctx, p.id, p.name, p.email, p.company, p.role); err != nil {
			return err
		}
	}

	// Relationships — mix of investment, advisor, colleague, intro types
	rels := []struct {
		from, to, relType string
		strength          int
	}{
		{"p1", "p2", "investor", 85},
		{"p1", "p4", "investor", 72},
		{"p1", "p5", "colleague", 90},
		{"p2", "p3", "advisor", 65},
		{"p2", "p6", "colleague", 95},
		{"p3", "p4", "friend", 80},
		{"p3", "p7", "former_colleague", 55},
		{"p4", "p8", "colleague", 92},
		{"p5", "p7", "intro", 40},
		{"p6", "p8", "conference", 30},
		{"p7", "p1", "advisor", 60},
	}
	for _, r := range rels {
		if err := c.CreateRelationship(ctx, r.from, r.to, r.relType, r.strength); err != nil {
			return err
		}
	}

	log.Println("Demo graph seeded successfully")
	return nil
}
