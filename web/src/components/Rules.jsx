import { useState, useEffect } from 'react'
import axios from 'axios'

export default function Rules() {
  const [rules, setRules] = useState([])
  const [showForm, setShowForm] = useState(false)
  const [loading, setLoading] = useState(true)
  const [newRule, setNewRule] = useState({
    domain: '',
    action: 'ALLOW',
    type: 'DOMAIN',
  })

  useEffect(() => {
    fetchRules()
  }, [])

  const fetchRules = async () => {
    try {
      const response = await axios.get('/api/rules')
      setRules(Array.isArray(response.data) ? response.data : [])
    } catch (err) {
      console.error('Failed to fetch rules')
    } finally {
      setLoading(false)
    }
  }

  const handleAddRule = async (e) => {
    e.preventDefault()
    try {
      await axios.post('/api/rules', newRule)
      setNewRule({ domain: '', action: 'ALLOW', type: 'DOMAIN' })
      setShowForm(false)
      fetchRules()
    } catch (err) {
      console.error('Failed to add rule')
    }
  }

  const handleDeleteRule = async (ruleId) => {
    try {
      await axios.delete(`/api/rules/${ruleId}`)
      fetchRules()
    } catch (err) {
      console.error('Failed to delete rule')
    }
  }

  if (loading) return <div className="text-center py-8">Loading rules...</div>

  return (
    <div className="space-y-4">
      <button
        onClick={() => setShowForm(!showForm)}
        className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium"
      >
        {showForm ? '✕ Cancel' : '+ Add Rule'}
      </button>

      {showForm && (
        <form onSubmit={handleAddRule} className="bg-white rounded-lg shadow p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Domain / Pattern</label>
            <input
              type="text"
              value={newRule.domain}
              onChange={(e) => setNewRule({ ...newRule, domain: e.target.value })}
              placeholder="e.g., ads.example.com"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Action</label>
            <select
              value={newRule.action}
              onChange={(e) => setNewRule({ ...newRule, action: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="ALLOW">Allow</option>
              <option value="BLOCK">Block</option>
              <option value="REDIRECT">Redirect</option>
            </select>
          </div>
          <button
            type="submit"
            className="w-full px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg font-medium"
          >
            Save Rule
          </button>
        </form>
      )}

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-100 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Pattern</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Type</th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Action</th>
              <th className="px-6 py-3 text-right text-sm font-semibold text-gray-900">Actions</th>
            </tr>
          </thead>
          <tbody>
            {rules.length > 0 ? (
              rules.map((rule) => (
                <tr key={rule.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="px-6 py-3 text-sm text-gray-900">{rule.domain || rule.pattern || 'N/A'}</td>
                  <td className="px-6 py-3 text-sm text-gray-600">{rule.type || 'DOMAIN'}</td>
                  <td className="px-6 py-3 text-sm">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-semibold ${
                      rule.action === 'BLOCK'
                        ? 'bg-red-100 text-red-800'
                        : rule.action === 'ALLOW'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }`}>
                      {rule.action}
                    </span>
                  </td>
                  <td className="px-6 py-3 text-right">
                    <button
                      onClick={() => handleDeleteRule(rule.id)}
                      className="text-red-600 hover:text-red-800 text-sm font-medium"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="4" className="px-6 py-4 text-center text-gray-500">No rules configured</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
