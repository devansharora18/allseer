import { useState, useEffect } from 'react'
import axios from 'axios'

export default function Dashboard() {
  const [stats, setStats] = useState(null)
  const [firewallEnabled, setFirewallEnabled] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetchStats()
    const interval = setInterval(fetchStats, 2000)
    return () => clearInterval(interval)
  }, [])

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/stats')
      setStats(response.data)
      setError(null)
    } catch (err) {
      setError('Failed to fetch stats')
    } finally {
      setLoading(false)
    }
  }

  const toggleFirewall = () => {
    setFirewallEnabled(!firewallEnabled)
  }

  if (loading) return <div className="text-center py-8">Loading...</div>

  return (
    <div className="space-y-6">
      {/* Status Card */}
      <div className="bg-white rounded-lg shadow p-6 border-l-4 border-blue-500">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Firewall Status</h2>
            <p className="text-sm text-gray-600 mt-1">
              {firewallEnabled ? '✅ Firewall Active' : '⚪ Firewall Inactive'}
            </p>
          </div>
          <button
            onClick={toggleFirewall}
            className={`px-6 py-2 rounded-lg font-medium transition-colors ${
              firewallEnabled
                ? 'bg-red-500 hover:bg-red-600 text-white'
                : 'bg-green-500 hover:bg-green-600 text-white'
            }`}
          >
            {firewallEnabled ? 'Disable' : 'Enable'}
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <StatCard label="Active Connections" value={stats.activeConnections || 0} icon="🔗" />
          <StatCard label="Blocked Requests" value={stats.blockedRequests || 0} icon="🚫" />
          <StatCard label="Allowed Requests" value={stats.allowedRequests || 0} icon="✅" />
          <StatCard label="Traffic (MB)" value={((stats.totalBytesForwarded || 0) / 1024 / 1024).toFixed(2)} icon="📊" />
        </div>
      )}

      {/* Activity */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-4">Recent Activity</h3>
        <div className="space-y-3">
          <ActivityRow domain="example.com" action="Allowed" timestamp="2 mins ago" />
          <ActivityRow domain="ads.tracker.net" action="Blocked" timestamp="5 mins ago" />
          <ActivityRow domain="github.com" action="Allowed" timestamp="10 mins ago" />
        </div>
      </div>

      {error && <div className="bg-red-50 border border-red-200 rounded p-4 text-red-800">{error}</div>}
    </div>
  )
}

function StatCard({ label, value, icon }) {
  return (
    <div className="bg-white rounded-lg shadow p-4 border-t-4 border-blue-500">
      <div className="text-3xl mb-2">{icon}</div>
      <div className="text-2xl font-bold text-gray-900">{value}</div>
      <div className="text-sm text-gray-600">{label}</div>
    </div>
  )
}

function ActivityRow({ domain, action, timestamp }) {
  const actionColor = action === 'Blocked' ? 'text-red-600' : 'text-green-600'
  return (
    <div className="flex items-center justify-between py-2 border-b border-gray-100">
      <div>
        <div className="font-semibold text-gray-900">{domain}</div>
        <div className="text-xs text-gray-500">{timestamp}</div>
      </div>
      <span className={`font-semibold ${actionColor}`}>{action}</span>
    </div>
  )
}
