'use client'

import { useEffect, useState } from 'react'
import { Plus, Loader2, AlertCircle } from 'lucide-react'

interface Person {
  id: string
  name: string
  title: string
  company: string
  email: string
}

export default function PeopleDirectory() {
  const [people, setPeople] = useState<Person[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [newPerson, setNewPerson] = useState({ name: '', title: '', company: '', email: '' })
  const [adding, setAdding] = useState(false)

  useEffect(() => {
    fetchPeople()
  }, [])

  const fetchPeople = async () => {
    try {
      setLoading(true)
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/people', {
        headers: { 'Authorization': `Bearer ${token}` },
      })
      if (!response.ok) throw new Error('Failed to fetch people')
      const data = await response.json()
      setPeople(data.people || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load people')
    } finally {
      setLoading(false)
    }
  }

  const handleAddPerson = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newPerson.name || !newPerson.email) return

    try {
      setAdding(true)
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/people', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newPerson),
      })
      if (!response.ok) throw new Error('Failed to add person')
      
      setNewPerson({ name: '', title: '', company: '', email: '' })
      await fetchPeople()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add person')
    } finally {
      setAdding(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="w-6 h-6 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Add Person Form */}
      <div className="bg-card border border-border rounded-lg p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
          <Plus className="w-5 h-5 text-primary" />
          Add New Person
        </h2>
        <form onSubmit={handleAddPerson} className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input
            type="text"
            placeholder="Full Name"
            value={newPerson.name}
            onChange={(e) => setNewPerson({ ...newPerson, name: e.target.value })}
            className="px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary"
            required
          />
          <input
            type="email"
            placeholder="Email"
            value={newPerson.email}
            onChange={(e) => setNewPerson({ ...newPerson, email: e.target.value })}
            className="px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary"
            required
          />
          <input
            type="text"
            placeholder="Job Title"
            value={newPerson.title}
            onChange={(e) => setNewPerson({ ...newPerson, title: e.target.value })}
            className="px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary"
          />
          <input
            type="text"
            placeholder="Company"
            value={newPerson.company}
            onChange={(e) => setNewPerson({ ...newPerson, company: e.target.value })}
            className="px-4 py-2 bg-input border border-border rounded-md text-foreground placeholder-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary"
          />
          <button
            type="submit"
            disabled={adding}
            className="md:col-span-2 px-4 py-2 bg-primary text-primary-foreground font-medium rounded-md hover:opacity-90 disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {adding && <Loader2 className="w-4 h-4 animate-spin" />}
            {adding ? 'Adding...' : 'Add Person'}
          </button>
        </form>
      </div>

      {/* Error Alert */}
      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <AlertCircle className="w-5 h-5 text-red-500" />
          <p className="text-sm text-red-500">{error}</p>
        </div>
      )}

      {/* People List */}
      <div className="space-y-3">
        <h2 className="text-lg font-semibold text-foreground">People in Network</h2>
        {people.length === 0 ? (
          <div className="text-center py-12 bg-card border border-border rounded-lg">
            <p className="text-muted-foreground">No people added yet</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {people.map((person) => (
              <div key={person.id} className="bg-card border border-border rounded-lg p-4 hover:border-primary transition">
                <h3 className="font-semibold text-foreground">{person.name}</h3>
                <p className="text-sm text-primary">{person.title}</p>
                <p className="text-sm text-muted-foreground">{person.company}</p>
                <p className="text-xs text-muted-foreground mt-2">{person.email}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
