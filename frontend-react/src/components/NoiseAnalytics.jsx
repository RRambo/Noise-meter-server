import React, { useState, useEffect } from 'react';
import { AreaChart, Area, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, LineChart, Line } from 'recharts';
import '../styles/NoiseAnalytics.css';
import { dataAPI } from '../services/api';
import ChartEmptyState from './ChartEmptyState';

function NoiseAnalytics({ roomName, allLocations }) {
    const [activeTab, setActiveTab] = useState('daily');
    const [selectedDate, setSelectedDate] = useState(new Date());
    const [selectedRoom, setSelectedRoom] = useState(roomName);
    const [weekOffset, setWeekOffset] = useState(0);
    const [chartData, setChartData] = useState([]);
    const [quietTimeData, setQuietTimeData] = useState([]);

    // Update selected room when roomName changes
    useEffect(() => {
        setSelectedRoom(roomName);
    }, [roomName]);

    // Get which day of week from date
    const getDayOfWeek = (date) => {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        return days[date.getDay()];
    };

    // Get which daye from chosen day of week
    const getDateForDay = (dayName) => {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        const targetDayIndex = days.indexOf(dayName);

        if (targetDayIndex === -1) return selectedDate;

        const today = new Date();
        const currentDayIndex = today.getDay();
        const diff = targetDayIndex - currentDayIndex;

        const newDate = new Date(today);
        newDate.setDate(today.getDate() + diff);

        return newDate;
    };

    // Which day of week is selected date
    const currentDayOfWeek = getDayOfWeek(selectedDate);

    // Handle date change
    const handleDateChange = (newDate) => {
        if (newDate == "") {
            return;
        } else {
            setSelectedDate(newDate);
        }
    };

    // Handle week selection change
    const handleDayChange = (dayName) => {
        const newDate = getDateForDay(dayName);
        setSelectedDate(newDate);
    };

    // Generate chart data
    useEffect(() => {
        generateChartData();
    }, [activeTab, selectedDate, selectedRoom, weekOffset]);



    const getDailySummary = async (room, date) => {
        // Gets data from a specific room measured during a specific day
        try {
            // all this to set the date to midnight without changing the selectedDate
            let y, m, d;
            y = date.getFullYear();
            m = date.getMonth() + 1;
            d = date.getDate();

            const utcMid = new Date(Date.UTC(y, m - 1, d, 0, 0, 0, 0));

            const response = await dataAPI.getDailySummary(room, utcMid.toISOString());
            // Do something with the fetched data
            return response.data
        } catch (error) {
            console.error('Error fetching daily summary:', error);
        }
    };

    const generateChartData = async () => {
        // too laggy but works :/
        // get data from a specific room with measure_time during a specific day
        const dailyData = await getDailySummary(selectedRoom, selectedDate) || []

        // Check if there are data within the hours of 8 and 17
        const hasDataInWindow = dailyData.some(e => {
            const h = new Date(e.measure_time).getHours();
            return h >= 8 && h <= 17;
        });

        if (!hasDataInWindow) {
            setChartData([])
            return;
        }

        if (activeTab === 'daily') {
            // calculates the average and looks for the peaks
            // Works while we wait for the hourly data implementation

            // Simulate daily data (8:00-17:00)
            const hours = Array.from({ length: 10 }, (_, i) => 8 + i);
            // buckets for data within each hour
            const buckets = Object.fromEntries(hours.map(h => [h, { count: 0, sum: 0, max: null }]));
            const data = [];
            // fill the buckets with data from the server
            dailyData.forEach(e => {
                // ignore timezone offset
                const match = e.measure_time.match(/T(\d{1,2}):/);
                const hour = match ? Number(match[1]) : NaN;
                
                // Only add the data within 8 - 17
                if (hours.includes(hour)) {
                    const b = buckets[hour];

                    b.count += 1;
                    b.sum += e.sound_level;
                    b.max = b.max === null ? e.sound_level : Math.max(b.max, e.sound_level);
                }
            });
            // Calculate the average and peak to the data array
            hours.map(h => {
                const { count, sum, max } = buckets[h];
                data.push({
                    time: `${h}:00`,
                    avgLevel: count > 0 ? Math.round(sum / count) : null,
                    peakLevel: max === null ? null : max,
                });
            });
            /*// for checking the data in the buckets
            data.forEach(e =>
                console.log(`time: ${e.time}, average: ${e.avgLevel}, peak: ${e.peakLevel}`)
            )*/
            setChartData(data);
        } else {
            // Generate weekly bar chart data
            const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];
            const data = days.map(day => {
                const avgNoise = Math.round(50 + Math.random() * 25);
                const peakNoise = Math.round(avgNoise + 15 + Math.random() * 20);
                return {
                    day: day,
                    avgNoise: avgNoise,
                    peakNoise: peakNoise
                };
            });
            setChartData(data);

            // Generate quiet time duration data
            const quietData = days.map(day => {
                const duration = Math.round(240 - Math.random() * 120);
                return {
                    day: day,
                    duration: duration
                };
            });
            setQuietTimeData(quietData);
        }

    };

    const getWeekDateRange = () => {
        const today = new Date();
        const currentDay = today.getDay(); // 0 = Sunday, 1 = Monday, etc.
        const monday = new Date(today);
        monday.setDate(today.getDate() - currentDay + (currentDay === 0 ? -6 : 1) + (weekOffset * 7));

        const friday = new Date(monday);
        friday.setDate(monday.getDate() + 4);

        const formatDate = (date) => {
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            return `${months[date.getMonth()]} ${date.getDate()}`;
        };

        return `${formatDate(monday)} - ${formatDate(friday)}, ${monday.getFullYear()}`;
    };

    const getWeekLabel = () => {
        if (weekOffset === 0) return 'Current Week';
        if (weekOffset === -1) return 'Last Week';
        return `${Math.abs(weekOffset)} Weeks Ago`;
    };

    const canGoBack = weekOffset > -4;
    const canGoForward = weekOffset < 0;

    const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'];

    return (
        <div className="noise-analytics-card card">
            <div className="card-body">
                {/* Header */}
                <div className="analytics-header">
                    <div>
                        <h5 className="card-title mb-1">
                            <span className="icon">📊</span> Noise Analysis
                        </h5>
                        <p className="card-subtitle text-muted">
                            Track noise patterns to plan activities during quieter periods
                        </p>
                    </div>
                    <div className="room-selector">
                        <label className="form-label mb-1">Room</label>
                        <select
                            className="form-select form-select-sm"
                            value={selectedRoom}
                            onChange={(e) => setSelectedRoom(e.target.value)}
                        >
                            {allLocations && allLocations.length > 0 ? (
                                allLocations.map(loc => (
                                    <option key={loc.id} value={loc.name}>{loc.name}</option>
                                ))
                            ) : (
                                <option>{roomName}</option>
                            )}
                        </select>
                    </div>
                </div>

                {/* Tab Buttons */}
                <div className="tab-buttons">
                    <button
                        className={`tab-btn ${activeTab === 'daily' ? 'active' : ''}`}
                        onClick={() => setActiveTab('daily')}
                    >
                        Daily Analysis
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'weekly' ? 'active' : ''}`}
                        onClick={() => setActiveTab('weekly')}
                    >
                        Weekly Analysis
                    </button>
                </div>

                {/* Selectors */}
                {activeTab === 'daily' ? (
                    <div className="selectors">
                        <div className="selector-item">
                            <label>Date</label>
                            <input
                                type="date"
                                className="form-control form-control-sm"
                                value={selectedDate.toISOString().split('T')[0]}
                                onChange={(e) => {
                                    // prevent error when clearing date
                                    e.target.value == "" ? handleDateChange(new Date()) : handleDateChange(new Date(e.target.value))
                                }}
                            />
                        </div>
                        <div className="selector-item">
                            <label>Choose Day from This Week</label>
                            <select
                                className="form-select form-select-sm"
                                value={currentDayOfWeek}
                                onChange={(e) => handleDayChange(e.target.value)}
                            >
                                {['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'].map(day => (
                                    <option key={day} value={day}>{day}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                ) : (
                    <div className="week-selector">
                        <div className="week-label">This Week</div>
                        <div className="week-navigation">
                            <button
                                className="nav-btn"
                                onClick={() => setWeekOffset(prev => prev - 1)}
                                disabled={!canGoBack}
                            >
                                &#8249;
                            </button>
                            <span className="week-text">{getWeekLabel()}</span>
                            <button
                                className="nav-btn"
                                onClick={() => setWeekOffset(prev => prev + 1)}
                                disabled={!canGoForward}
                            >
                                &#8250;
                            </button>
                        </div>
                        <div className="week-range">{getWeekDateRange()}</div>
                    </div>
                )}

                {/* Chart */}
                <div className="chart-container">
                    {activeTab === 'daily' ? (
                        chartData.length == 0 ? (
                            // Overkill of an empty chart state XD
                            <ChartEmptyState chosenDate={selectedDate.toLocaleDateString()} onRetry={() => handleDateChange(new Date(selectedDate))} />
                        ) : (
                            <ResponsiveContainer width="100%" height={300}>
                                <AreaChart data={chartData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                                    <defs>
                                        <linearGradient id="colorAvg" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#5F9EA0" stopOpacity={0.3} />
                                            <stop offset="95%" stopColor="#5F9EA0" stopOpacity={0.05} />
                                        </linearGradient>
                                    </defs>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                    <XAxis
                                        dataKey="time"
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                    />
                                    <YAxis
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                        domain={[0, 100]}
                                        label={{ value: 'Decibels (dB)', angle: -90, position: 'insideLeft', style: { fontSize: '0.85rem' } }}
                                    />
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: 'white',
                                            border: '1px solid #ccc',
                                            borderRadius: '8px',
                                            padding: '10px'
                                        }}
                                    />
                                    <Area
                                        type="monotone"
                                        dataKey="avgLevel"
                                        stroke="#5F9EA0"
                                        strokeWidth={2}
                                        fill="url(#colorAvg)"
                                        name="Average Level"
                                    />
                                    <Area
                                        type="monotone"
                                        dataKey="peakLevel"
                                        stroke="#5F9EA0"
                                        strokeWidth={2}
                                        strokeDasharray="5 5"
                                        fill="none"
                                        name="Peak Level"
                                    />
                                </AreaChart>
                            </ResponsiveContainer>
                        )

                    ) : (
                        <>
                            <ResponsiveContainer width="100%" height={300}>
                                <BarChart data={chartData} margin={{ top: 20, right: 30, left: 0, bottom: 5 }}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                    <XAxis
                                        dataKey="day"
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                    />
                                    <YAxis
                                        stroke="#666"
                                        style={{ fontSize: '0.85rem' }}
                                        domain={[0, 100]}
                                        label={{ value: 'Decibels (dB)', angle: -90, position: 'insideLeft', style: { fontSize: '0.85rem' } }}
                                    />
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: 'white',
                                            border: '1px solid #ccc',
                                            borderRadius: '8px',
                                            padding: '10px'
                                        }}
                                    />
                                    <Legend
                                        wrapperStyle={{ paddingTop: '10px' }}
                                        iconType="rect"
                                    />
                                    <Bar dataKey="avgNoise" fill="#81C9CC" name="Average Noise" />
                                    <Bar dataKey="peakNoise" fill="#5F9EA0" name="Peak Noise" />
                                </BarChart>
                            </ResponsiveContainer>

                            {/* Quiet Time Duration Chart */}
                            <div className="quiet-time-section">
                                <h6 className="section-title">
                                    <span className="icon">🔇</span> Quiet Time Duration (minutes)
                                </h6>
                                <ResponsiveContainer width="100%" height={200}>
                                    <LineChart data={quietTimeData} margin={{ top: 10, right: 30, left: 0, bottom: 5 }}>
                                        <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                                        <XAxis
                                            dataKey="day"
                                            stroke="#666"
                                            style={{ fontSize: '0.85rem' }}
                                        />
                                        <YAxis
                                            stroke="#666"
                                            style={{ fontSize: '0.85rem' }}
                                            domain={[0, 300]}
                                        />
                                        <Tooltip
                                            contentStyle={{
                                                backgroundColor: 'white',
                                                border: '1px solid #ccc',
                                                borderRadius: '8px',
                                                padding: '10px'
                                            }}
                                        />
                                        <Line
                                            type="monotone"
                                            dataKey="duration"
                                            stroke="#5F9EA0"
                                            strokeWidth={2}
                                            dot={{ fill: '#5F9EA0', r: 4 }}
                                        />
                                    </LineChart>
                                </ResponsiveContainer>
                            </div>
                        </>
                    )}
                </div>

                {/* Legend for Daily */}
                {activeTab === 'daily' && (
                    <div className="chart-legend">
                        <div className="legend-item">
                            <span className="legend-line solid"></span>
                            <span>Average Level</span>
                        </div>
                        <div className="legend-item">
                            <span className="legend-line dashed"></span>
                            <span>Peak Level</span>
                        </div>
                    </div>
                )}

                {/* Planning Tip */}
                <div className="planning-tip">
                    <strong>💡 Planning Tip:</strong> {activeTab === 'daily'
                        ? 'Schedule quiet activities (story time, nap time) during naturally quieter periods (8:00-9:00, 12:00-13:00) and active play during peak energy times (10:00-11:00).'
                        : 'Wednesday shows the longest quiet periods and lowest average noise - ideal for introducing new concepts or activities requiring focus. Friday has higher energy levels - perfect for group activities and celebrations.'}
                </div>
            </div>
        </div>
    );
}

export default NoiseAnalytics;