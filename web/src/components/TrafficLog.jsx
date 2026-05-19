import { useState, useEffect } from 'react'
import axios from 'axios'

export default function TrafficLog() {
  const [logs, setLogs] = useState([])
  const [filter, setFilter] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchLogs()
    const interval = setInterval(fetchLogs, 3000)
    return () => clearInterval(interval)
  }, [])

  const fetchLogs = async () => {
    try {
      const response = await axios.get('/admin/logs?limit=50')
      setLogs(Array.isArray(response.data) ? response.data : [])
    } catch (err) {
      console.error('Failed to fetch logs')
    } finally {
      setLoading(false)
    }
  }

  const filteredLogs = logs.filter((log) =>
    log.domain?.toLowerCase().includes(filter.toLowerCase()) ||
    log.source?.toLowerCase().includes(filter.toLowerCase())
  )

  if (loading) return <div className="text-center py-8">Loading traffic log...</div>

  return (
    <div className="space-y-4">
      <div className="bg-white rounded-lg shadow p-4">
        <input
          type="text"
          placeholder="Filter by domain or IP..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-100 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Domain</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Source IP</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Action</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Time</th>
            </tr>
          </thead>
          <tbody>
            {filteredLogs.length > 0 ? (
              filteredLogs.map((log, idx) => (
                <tr key={idx} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="px-6 py-3 text-sm text-gray-900 font-mono">{log.domain || 'N/A'}</td>
                  <td className="px-6 py-3 text-sm text-gray-600 font-mono">{log.source || 'N/A'}</td>
                  <td className="px-6 py-3 text-sm">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-semibold ${
                      log.action === 'BLOCKED' 
                        ? 'bg-red-100 text-red-800' 
                        : 'bg-green-100 text-green-800'
                    }`}>
                      {log.action || 'UNKNOWN'}
                    </span>
                  </td>
                  <td className="px-6 py-3 text-sm text-gray-500">{new Date(log.timestamp).toLocaleTimeString()}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="4" className="px-6 py-4 text-center text-gray-500">No traffic logs found</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
