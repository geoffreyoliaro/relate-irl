'use client'

import { useEffect, useState } from 'react'
import { Network, Loader2, AlertCircle } from 'lucide-react'

export default function NetworkVisualization() {
  const [loading, setLoading] = useState(false)
  const [selectedPerson, setSelectedPerson] = useState('')
  const [networkData, setNetworkData] = useState<any>(null)
  const [error, setError] = useState('')
  const [people, setPeople] = useState<any[]>([])

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

  const fetchNetwork = async () => {
    if (!selectedPerson) return

    try {
      setLoading(true)
      setError('')
      const token = localStorage.getItem('token')
      const response = await fetch(`http://localhost:8080/api/v1/relationships/${selectedPerson}/network`, {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (!response.ok) throw new Error('Failed to fetch network')
      const data = await response.json()
      setNetworkData(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load network')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Network Controls */}
      <div className="bg-card border border-border rounded-lg p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
          <Network className="w-5 h-5 text-primary" />
          Explore Network
        </h2>
        <div className="flex gap-4 flex-wrap">
          <select
            value={selectedPerson}
            onChange={(e) => setSelectedPerson(e.target.value)}
            className="flex-1 min-w-64 px-4 py-2 bg-input border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">Select a person</option>
            {people.map((person) => (
              <option key={person.id} value={person.id}>
                {person.name}
              </option>
            ))}
          </select>
          <button
            onClick={fetchNetwork}
            disabled={!selectedPerson || loading}
            className="px-6 py-2 bg-primary text-primary-foreground font-medium rounded-md hover:opacity-90 disabled:opacity-50 flex items-center gap-2"
          >
            {loading && <Loader2 className="w-4 h-4 animate-spin" />}
            {loading ? 'Loading...' : 'View Network'}
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

      {/* Network Visualization */}
      {networkData && (
        <div className="space-y-6">
          {/* Center Node */}
          <div className="bg-card border border-primary rounded-lg p-6">
            <h3 className="text-sm font-semibold text-muted-foreground mb-3">Center Node</h3>
            <div className="bg-gradient-to-br from-primary/20 to-primary/5 rounded-lg p-4 border border-primary/20">
              <p className="font-semibold text-foreground text-lg">{networkData.center?.name}</p>
              <p className="text-sm text-primary">{networkData.center?.title}</p>
            </div>
          </div>

          {/* Connections */}
          <div>
            <h3 className="text-sm font-semibold text-muted-foreground mb-3">
              Connected Nodes ({networkData.connections?.length || 0})
            </h3>
            {networkData.connections && networkData.connections.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {networkData.connections.map((conn: any, idx: number) => (
                  <div key={idx} className="bg-card border border-border rounded-lg p-4 hover:border-primary/50 transition">
                    <div className="flex items-start justify-between mb-2">
                      <h4 className="font-medium text-foreground">{conn.name}</h4>
                      <span className="text-xs bg-primary/20 text-primary px-2 py-1 rounded">
                        {conn.distance} hops
                      </span>
                    </div>
                    <p className="text-sm text-muted-foreground">Via: {conn.via}</p>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 bg-card border border-border rounded-lg">
                <p className="text-muted-foreground">No direct connections found</p>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Empty State */}
      {!networkData && !loading && (
        <div className="text-center py-16 bg-card border border-border rounded-lg">
          <Network className="w-12 h-12 text-muted-foreground mx-auto mb-4 opacity-50" />
          <p className="text-muted-foreground">Select a person to visualize their network</p>
        </div>
      )}
    </div>
  )
}
