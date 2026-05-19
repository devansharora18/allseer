import { useState } from 'react'
import Dashboard from './components/Dashboard'
import TrafficLog from './components/TrafficLog'
import Rules from './components/Rules'
import Settings from './components/Settings'

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard')

  const tabs = [
    { id: 'dashboard', label: 'Dashboard', icon: '📊' },
    { id: 'traffic', label: 'Traffic Log', icon: '🔍' },
    { id: 'rules', label: 'Rules', icon: '⚙️' },
    { id: 'settings', label: 'Settings', icon: '🔧' },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <h1 className="text-3xl font-bold text-gray-900">🔥 Allseer Firewall</h1>
          <p className="text-gray-600 mt-1">Local programmable firewall & proxy</p>
        </div>
      </header>

      {/* Navigation Tabs */}
      <nav className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 flex gap-1 overflow-x-auto">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300'
              }`}
            >
              {tab.icon} {tab.label}
            </button>
          ))}
        </div>
      </nav>

      {/* Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        {activeTab === 'dashboard' && <Dashboard />}
        {activeTab === 'traffic' && <TrafficLog />}
        {activeTab === 'rules' && <Rules />}
        {activeTab === 'settings' && <Settings />}
      </main>
    </div>
  )
}
