'use client'

import { useEffect, useState } from 'react'
import { Zap, Loader2, AlertCircle } from 'lucide-react'

export default function IntelligenceQueries() {
  const [people, setPeople] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  
  const [queryType, setQueryType] = useState('mutual')
  const [personA, setPersonA] = useState('')
  const [personB, setPersonB] = useState('')
  const [results, setResults] = useState<any>(null)

  useEffect(() => {
    fetchPeople()
  }, [])

  const fetchPeople = async () => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/people', {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (!response.ok) throw new Error('Failed to fetch people')
      const data = await response.json()
      setPeople(data.people || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load people')
    }
  }

  const runQuery = async () => {
    if (!personA || !personB) return

    try {
      setLoading(true)
      setError('')
      const token = localStorage.getItem('token')
      
      let endpoint = ''
      if (queryType === 'mutual') {
        endpoint = `http://localhost:8080/api/v1/intelligence/mutual-connections?person_a=${personA}&person_b=${personB}`
      } else if (queryType === 'path') {
        endpoint = `http://localhost:8080/api/v1/intelligence/shortest-path?from=${personA}&to=${personB}`
      } else {
        endpoint = `http://localhost:8080/api/v1/intelligence/relationship-strength?from=${personA}&to=${personB}`
      }

      const response = await fetch(endpoint, {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      
      if (!response.ok) throw new Error('Query failed')
      const data = await response.json()
      setResults(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Query failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Query Builder */}
      <div className="bg-card border border-border rounded-lg p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
          <Zap className="w-5 h-5 text-primary" />
          Intelligence Queries
        </h2>

        <div className="space-y-4">
          {/* Query Type */}
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">Query Type</label>
            <select
              value={queryType}
              onChange={(e) => {
                setQueryType(e.target.value)
                setResults(null)
              }}
              className="w-full px-4 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
            >
              <option value="mutual">Mutual Connections</option>
              <option value="path">Shortest Path</option>
              <option value="strength">Relationship Strength</option>
            </select>
          </div>

          {/* Person Selection */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">Person A</label>
              <select
                value={personA}
                onChange={(e) => setPersonA(e.target.value)}
                className="w-full px-4 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
              >
                <option value="">Select person</option>
                {people.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">Person B</label>
              <select
                value={personB}
                onChange={(e) => setPersonB(e.target.value)}
                className="w-full px-4 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
              >
                <option value="">Select person</option>
                {people.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {/* Run Button */}
          <button
            onClick={runQuery}
            disabled={!personA || !personB || loading}
            className="w-full px-4 py-2 bg-primary text-primary-foreground font-medium rounded-md hover:opacity-90 disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {loading && <Loader2 className="w-4 h-4 animate-spin" />}
            {loading ? 'Running Query...' : 'Run Query'}
          </button>
        </div>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <AlertCircle className="w-5 h-5 text-red-500" />
          <p className="text-sm text-red-500">{error}</p>
        </div>
      )}

      {/* Results */}
      {results && (
        <div className="bg-card border border-border rounded-lg p-6">
          <h3 className="text-lg font-semibold text-foreground mb-4">Query Results</h3>

          {queryType === 'mutual' && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="bg-secondary/50 rounded p-3">
                  <p className="text-xs text-muted-foreground">Person A</p>
                  <p className="font-medium text-foreground">
                    {people.find((p) => p.id === results.person_a)?.name}
                  </p>
                </div>
                <div className="bg-secondary/50 rounded p-3">
                  <p className="text-xs text-muted-foreground">Person B</p>
                  <p className="font-medium text-foreground">
                    {people.find((p) => p.id === results.person_b)?.name}
                  </p>
                </div>
              </div>
              <div>
                <p className="text-sm font-medium text-foreground mb-2">
                  Mutual Connections: <span className="text-primary">{results.count}</span>
                </p>
                {results.mutual_connections?.length > 0 && (
                  <div className="space-y-2">
                    {results.mutual_connections.map((conn: any, idx: number) => (
                      <div key={idx} className="bg-secondary/30 rounded p-3">
                        <p className="font-medium text-foreground">{conn.name}</p>
                        <p className="text-xs text-muted-foreground">Path: {conn.connection_path?.join(' → ')}</p>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}

          {queryType === 'path' && (
            <div className="space-y-4">
              <div className="flex items-center gap-2 mb-4">
                <span className="text-sm font-medium text-muted-foreground">Path Length:</span>
                <span className="text-lg font-bold text-primary">{results.length}</span>
              </div>
              <div className="space-y-2">
                {results.path?.map((person: any, idx: number) => (
                  <div key={idx} className="flex items-center gap-2">
                    <div className="bg-primary/20 rounded-full w-8 h-8 flex items-center justify-center text-xs font-medium text-primary">
                      {idx + 1}
                    </div>
                    <p className="font-medium text-foreground">{person.name}</p>
                    {idx < results.path.length - 1 && (
                      <span className="text-muted-foreground ml-auto">→</span>
                    )}
                  </div>
                ))}
              </div>
              <div className="bg-secondary/30 rounded p-3 mt-4">
                <p className="text-xs text-muted-foreground">Overall Strength</p>
                <p className="text-lg font-bold text-primary">{(results.strength * 100).toFixed(1)}%</p>
              </div>
            </div>
          )}

          {queryType === 'strength' && (
            <div className="space-y-4">
              <div className="bg-gradient-to-r from-primary/20 to-accent/20 rounded-lg p-6">
                <p className="text-sm text-muted-foreground mb-2">Relationship Strength</p>
                <p className="text-4xl font-bold text-primary mb-2">{(results.strength * 100).toFixed(1)}%</p>
                <div className="w-full bg-secondary rounded-full h-2">
                  <div
                    className="bg-primary h-2 rounded-full transition-all"
                    style={{ width: `${results.strength * 100}%` }}
                  ></div>
                </div>
              </div>
              {results.factors && (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                  {Object.entries(results.factors).map(([key, value]: [string, any]) => (
                    <div key={key} className="bg-secondary/30 rounded p-3">
                      <p className="text-xs text-muted-foreground capitalize">{key.replace(/_/g, ' ')}</p>
                      <p className="text-lg font-bold text-primary">{(value * 100).toFixed(0)}%</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Empty State */}
      {!results && !loading && (
        <div className="text-center py-16 bg-card border border-border rounded-lg">
          <Zap className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-50" />
          <p className="text-muted-foreground">Select two people to run an intelligence query</p>
        </div>
      )}
    </div>
  )
}
