import React from 'react';

function StatsCards({ dailyPeak, weeklyAverage, monitoringRoom, isActive }) {
  return (
    <div className="row g-4 mb-4">
      {/* Daily Peak Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📈</span>
              <span className="stat-label">Daily Peak</span>
            </div>
            <div className="stat-value">{dailyPeak} dB</div>
            <div className="stat-description">Highest level today</div>
          </div>
        </div>
      </div>

      {/* Weekly Average Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📊</span>
              <span className="stat-label">Weekly Average</span>
            </div>
            <div className="stat-value">{weeklyAverage} dB</div>
            <div className="stat-description">This week's average</div>
          </div>
        </div>
      </div>

      {/* Monitoring Status Card */}
      <div className="col-md-4">
        <div className="stat-card card h-100">
          <div className="card-body">
            <div className="stat-header">
              <span className="stat-icon">📡</span>
              <span className="stat-label">Monitoring</span>
            </div>
            <div className="stat-value-room">{monitoringRoom}</div>
            {/* <div className="stat-description">
              <span 
                className={`status-indicator ${isActive ? 'active' : 'inactive'}`}
              >
                {isActive ? '● Currently active' : '○ Inactive'}
              </span>
            </div> */}
          </div>
        </div>
      </div>

      <style jsx>{`
        .stat-card {
          border-radius: 12px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
          border: none;
          transition: transform 0.2s, box-shadow 0.2s;
        }
      
        .card-body {
          padding: 1.5rem;
        }

        .stat-header {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          margin-bottom: 1rem;
        }

        .stat-icon {
          font-size: 1.5rem;
        }

        .stat-label {
          font-size: 0.9rem;
          color: #666;
          font-weight: 500;
        }

        .stat-value {
          font-size: 2.5rem;
          font-weight: 700;
          color: #5F9EA0;
          line-height: 1;
          margin-bottom: 0.5rem;
        }

        .stat-value-room {
          font-size: 1.5rem;
          font-weight: 600;
          color: #333;
          margin-bottom: 0.5rem;
          min-height: 2.5rem;
          display: flex;
          align-items: center;
        }

        .stat-description {
          font-size: 0.85rem;
          color: #888;
        }

        .status-indicator {
          display: inline-flex;
          align-items: center;
          gap: 0.25rem;
          font-weight: 500;
        }

        .status-indicator.active {
          color: #4CAF50;
        }

        .status-indicator.inactive {
          color: #999;
        }

        @media (max-width: 768px) {
          .stat-value {
            font-size: 2rem;
          }

          .stat-value-room {
            font-size: 1.25rem;
          }
        }
      `}</style>
    </div>
  );
}

export default StatsCards;