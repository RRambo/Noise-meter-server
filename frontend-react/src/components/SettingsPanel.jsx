import React from 'react';

function SettingsPanel({ 
  locations, 
  chosenLocation, 
  onLocationChange, 
  threshold, 
  onThresholdChange 
}) {
  return (
    <div className="settings-panel card">
      <div className="card-body">
        {/* Header */}
        <div className="settings-header mb-3">
          <h5 className="card-title">
            <span className="settings-icon">⚙️</span> Settings
          </h5>
          <p className="card-subtitle text-muted">
            Configure monitoring preferences
          </p>
        </div>

        {/* Listening Room Selector */}
        <div className="mb-4">
          <label className="form-label fw-semibold">Listening Room</label>
          <select 
            className="form-select"
            value={chosenLocation?.id || ''}
            onChange={(e) => onLocationChange(parseInt(e.target.value))}
          >
            {locations.length === 0 ? (
              <option value="">No locations available</option>
            ) : (
              locations.map(location => (
                <option key={location.id} value={location.id}>
                  {location.name}
                </option>
              ))
            )}
          </select>
        </div>

        {/* Noise Threshold Slider */}
        <div className="mb-3">
          <div className="d-flex justify-content-between align-items-center mb-2">
            <label className="form-label fw-semibold mb-0">Noise Threshold</label>
            <span className="badge bg-light text-dark">{threshold} dB</span>
          </div>
          
          <input 
            type="range" 
            className="form-range" 
            min="0" 
            max="120" 
            step="5"
            value={threshold}
            onChange={(e) => onThresholdChange(parseInt(e.target.value))}
          />
          
          <p className="text-muted small mt-2 mb-0">
            You'll be notified when noise exceeds this level
          </p>
        </div>
      </div>

      <style jsx>{`
        .settings-panel {
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          border: none;
        }

        .settings-icon {
          font-size: 1.2rem;
          margin-right: 4px;
        }

        .card-title {
          margin-bottom: 4px;
          font-size: 1.1rem;
          font-weight: 600;
        }

        .card-subtitle {
          font-size: 0.875rem;
        }

        .form-range {
          cursor: pointer;
        }

        .form-range::-webkit-slider-thumb {
          background-color: #5F9EA0;
        }

        .form-range::-moz-range-thumb {
          background-color: #5F9EA0;
        }

        .form-select:focus,
        .form-range:focus {
          border-color: #5F9EA0;
          box-shadow: 0 0 0 0.2rem rgba(95, 158, 160, 0.25);
        }

        .badge {
          font-size: 0.9rem;
          padding: 0.4em 0.8em;
        }
      `}</style>
    </div>
  );
}

export default SettingsPanel;