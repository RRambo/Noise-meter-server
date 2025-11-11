import React from 'react';

function NoiseMeter({ currentLevel, threshold, roomName }) {
  // Calculate percentage for the circular progress
  const percentage = Math.min((currentLevel / 120) * 100, 100);
  
  // Determine status based on level vs threshold
  const getStatus = () => {
    if (currentLevel < threshold * 0.7) {
      return { label: 'Quiet', color: '#5F9EA0', bgColor: '#E0F2F1' };
    } else if (currentLevel < threshold) {
      return { label: 'Moderate', color: '#FFA726', bgColor: '#FFF3E0' };
    } else {
      return { label: 'Loud', color: '#EF5350', bgColor: '#FFEBEE' };
    }
  };

  const status = getStatus();

  // SVG circle parameters
  const size = 280;
  const strokeWidth = 20;
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  return (
    <div className="noise-meter-card card">
      <div className="card-body">
        {/* Header */}
        <div className="d-flex justify-content-between align-items-center mb-4">
          <h5 className="card-title mb-0">Current Noise Level</h5>
          <span className="room-badge">{roomName}</span>
        </div>

        {/* Circular Meter */}
        <div className="meter-container">
          <svg width={size} height={size} className="circular-meter">
            {/* Background circle */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke="#e0e0e0"
              strokeWidth={strokeWidth}
            />
            
            {/* Progress circle */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke={status.color}
              strokeWidth={strokeWidth}
              strokeDasharray={circumference}
              strokeDashoffset={strokeDashoffset}
              strokeLinecap="round"
              transform={`rotate(-90 ${size / 2} ${size / 2})`}
              style={{ transition: 'stroke-dashoffset 0.5s ease' }}
            />
          </svg>

          {/* Center content */}
          <div className="meter-content">
            <div className="sound-icon">🔊</div>
            <div className="noise-level">{currentLevel}</div>
            <div className="noise-unit">dB</div>
          </div>

          {/* Status badge below meter */}
          <div className="status-badge-container">
            <span 
              className="status-badge"
              style={{ 
                backgroundColor: status.bgColor, 
                color: status.color 
              }}
            >
              {status.label}
            </span>
          </div>
        </div>

        {/* Info Row */}
        <div className="info-row">
          <div className="info-item">
            <div className="info-label">Threshold</div>
            <div className="info-value">{threshold} dB</div>
          </div>
          <div className="info-divider"></div>
          <div className="info-item">
            <div className="info-label">Status</div>
            <div className="info-value">{status.label}</div>
          </div>
          <div className="info-divider"></div>
          <div className="info-item">
            <div className="info-label">Level</div>
            <div className="info-value">{currentLevel} dB</div>
          </div>
        </div>
      </div>

      <style jsx>{`
        .noise-meter-card {
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          border: none;
        }

        .room-badge {
          background-color: #f5f5f5;
          padding: 0.4rem 1rem;
          border-radius: 8px;
          font-size: 0.9rem;
          font-weight: 500;
        }

        .meter-container {
          position: relative;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          margin: 2rem 0;
        }

        .circular-meter {
          display: block;
        }

        .meter-content {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          text-align: center;
        }

        .sound-icon {
          font-size: 2.5rem;
          margin-bottom: 0.5rem;
          opacity: 0.6;
        }

        .noise-level {
          font-size: 4rem;
          font-weight: 700;
          line-height: 1;
          color: #333;
        }

        .noise-unit {
          font-size: 1.2rem;
          color: #666;
          margin-top: 0.25rem;
        }

        .status-badge-container {
          margin-top: 1.5rem;
        }

        .status-badge {
          padding: 0.5rem 1.5rem;
          border-radius: 20px;
          font-size: 0.9rem;
          font-weight: 500;
        }

        .info-row {
          display: flex;
          justify-content: space-around;
          align-items: center;
          padding: 1.5rem 0;
          border-top: 1px solid #e0e0e0;
          margin-top: 1rem;
        }

        .info-item {
          text-align: center;
          flex: 1;
        }

        .info-label {
          font-size: 0.85rem;
          color: #666;
          margin-bottom: 0.5rem;
        }

        .info-value {
          font-size: 1.1rem;
          font-weight: 600;
          color: #333;
        }

        .info-divider {
          width: 1px;
          height: 40px;
          background-color: #e0e0e0;
        }

        @media (max-width: 768px) {
          .noise-level {
            font-size: 3rem;
          }

          .info-row {
            flex-direction: column;
            gap: 1rem;
          }

          .info-divider {
            display: none;
          }
        }
      `}</style>
    </div>
  );
}

export default NoiseMeter;