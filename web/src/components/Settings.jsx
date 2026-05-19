import { useState } from 'react'

export default function Settings() {
  const [mitmEnabled, setMitmEnabled] = useState(false)
  const [loggingLevel, setLoggingLevel] = useState('INFO')

  const downloadCA = () => {
    alert('CA certificate download not yet implemented. Use /api/ca endpoint.')
  }

  return (
    <div className="space-y-6">
      {/* MITM Settings */}
      <div className="bg-white rounded-lg shadow p-6 border-l-4 border-yellow-500">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-lg font-bold text-gray-900">HTTPS Interception (MITM)</h3>
            <p className="text-sm text-gray-600 mt-1">
              Enable to intercept and inspect HTTPS traffic (requires CA trust)
            </p>
          </div>
          <button
            onClick={() => setMitmEnabled(!mitmEnabled)}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${
              mitmEnabled
                ? 'bg-yellow-500 hover:bg-yellow-600 text-white'
                : 'bg-gray-300 hover:bg-gray-400 text-white'
            }`}
          >
            {mitmEnabled ? 'Enabled' : 'Disabled'}
          </button>
        </div>
        {mitmEnabled && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded">
            <p className="text-sm text-yellow-800 mb-4">
              ⚠️ HTTPS MITM requires clients to trust the Allseer CA certificate. Download and install below.
            </p>
            <button
              onClick={downloadCA}
              className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded font-medium"
            >
              Download CA Certificate
            </button>
          </div>
        )}
      </div>

      {/* Logging Settings */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-bold text-gray-900 mb-4">Logging Level</h3>
        <select
          value={loggingLevel}
          onChange={(e) => setLoggingLevel(e.target.value)}
          className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="DEBUG">Debug (verbose)</option>
          <option value="INFO">Info (default)</option>
          <option value="WARN">Warning (errors only)</option>
          <option value="ERROR">Error (critical only)</option>
        </select>
        <p className="text-sm text-gray-600 mt-2">Current level: <span className="font-semibold">{loggingLevel}</span></p>
      </div>

      {/* Network Mode */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-bold text-gray-900 mb-4">Network Mode</h3>
        <div className="space-y-3">
          <label className="flex items-center">
            <input type="radio" name="mode" value="host" defaultChecked className="mr-3" />
            <div>
              <div className="font-semibold text-gray-900">Host Mode (Local)</div>
              <div className="text-sm text-gray-600">Filter traffic from this host only</div>
            </div>
          </label>
          <label className="flex items-center">
            <input type="radio" name="mode" value="gateway" className="mr-3" />
            <div>
              <div className="font-semibold text-gray-900">Gateway Mode (Network)</div>
              <div className="text-sm text-gray-600">Filter traffic from all connected devices</div>
            </div>
          </label>
        </div>
      </div>

      {/* System Info */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-bold text-gray-900 mb-4">System Information</h3>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between border-b pb-2">
            <span className="text-gray-600">Allseer Version</span>
            <span className="font-mono text-gray-900">0.1.0</span>
          </div>
          <div className="flex justify-between border-b pb-2">
            <span className="text-gray-600">Listen Address</span>
            <span className="font-mono text-gray-900">127.0.0.1:8080</span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Database</span>
            <span className="font-mono text-gray-900">SQLite (local)</span>
          </div>
        </div>
      </div>
    </div>
  )
}
