'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { LogOut, Users, Network, Zap } from 'lucide-react'
import PeopleDirectory from '@/components/PeopleDirectory'
import NetworkVisualization from '@/components/NetworkVisualization'
import IntelligenceQueries from '@/components/IntelligenceQueries'

type TabType = 'people' | 'network' | 'intelligence'

export default function DashboardPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState<TabType>('people')
  const [user, setUser] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')
    
    if (!token || !userData) {
      router.push('/')
      return
    }

    setUser(JSON.parse(userData))
    setLoading(false)
  }, [router])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/')
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-2 border-primary border-t-transparent"></div>
          <p className="text-muted-foreground mt-4">Loading...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-foreground">Relate IRL</h1>
            <p className="text-xs text-muted-foreground">Relationship Intelligence</p>
          </div>

          <div className="flex items-center gap-4">
            {user && (
              <div className="text-right mr-4">
                <p className="text-sm font-medium text-foreground">{user.name || user.email}</p>
                <p className="text-xs text-muted-foreground">{user.email}</p>
              </div>
            )}
            <button
              onClick={handleLogout}
              className="px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:bg-secondary hover:text-foreground transition flex items-center gap-2"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation Tabs */}
      <div className="border-b border-border bg-card/30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex gap-8">
            <button
              onClick={() => setActiveTab('people')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition flex items-center gap-2 ${
                activeTab === 'people'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              <Users className="w-4 h-4" />
              People Directory
            </button>
            <button
              onClick={() => setActiveTab('network')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition flex items-center gap-2 ${
                activeTab === 'network'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              <Network className="w-4 h-4" />
              Network Graph
            </button>
            <button
              onClick={() => setActiveTab('intelligence')}
              className={`py-4 px-1 border-b-2 font-medium text-sm transition flex items-center gap-2 ${
                activeTab === 'intelligence'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-muted-foreground hover:text-foreground'
              }`}
            >
              <Zap className="w-4 h-4" />
              Intelligence Queries
            </button>
          </div>
        </div>
      </div>

      {/* Content Area */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'people' && <PeopleDirectory />}
        {activeTab === 'network' && <NetworkVisualization />}
        {activeTab === 'intelligence' && <IntelligenceQueries />}
      </main>
    </div>
  )
}
