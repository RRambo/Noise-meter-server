import React, { useState, useEffect } from 'react';
import { dataAPI, locationAPI } from './services/api';
import SettingsPanel from './components/SettingsPanel';
import NoiseMeter from './components/NoiseMeter';
import StatsCards from './components/StatsCards';
import NoiseAnalytics from './components/NoiseAnalytics'; 
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

function App() {

  // State for locations
  const [locations, setLocations] = useState([]);
  const [chosenLocation, setChosenLocation] = useState(null);
  const [loading, setLoading] = useState(true);
  
  // State for settings - read value from localStorage
  const [threshold, setThreshold] = useState(() => {
    const savedThreshold = localStorage.getItem('noiseThreshold');
    return savedThreshold ? parseInt(savedThreshold) : 75;
  });

  // Placeholder state for current noise level
  // ! Need to replace with real data fetching from Arduino later
  const [currentNoiseLevel, setCurrentNoiseLevel] = useState(25);

  // Stats data state - use sessionStorage for dailyPeak
  const [dailyPeak, setDailyPeak] = useState(() => {
    const saved = sessionStorage.getItem('dailyPeak');
    return saved ? parseInt(saved) : 0;
  });
  const [weeklyAverage, setWeeklyAverage] = useState(63);

  // Save dailyPeak to sessionStorage 
  useEffect(() => {
    sessionStorage.setItem('dailyPeak', dailyPeak);
  }, [dailyPeak]);

  // Simulate current noise level updates every few seconds
  useEffect(() => {
    const interval = setInterval(() => {
      const randomLevel = Math.floor(Math.random() * 100); // 0-99 dB
      setCurrentNoiseLevel(randomLevel);
      
      // Update daily peak
      setDailyPeak(prev => Math.max(prev, randomLevel));
    }, 3000); // Updtate every 3 seconds

    return () => clearInterval(interval);
  }, []); 

  // Simulate noise level updates every few seconds (for testing)
  useEffect(() => {
    const interval = setInterval(() => {
      // Simulate weekly avarage change every 10 seconds
      setWeeklyAverage(prev => {
        const change = Math.floor(Math.random() * 5) - 2; // -2 to +2
        const newValue = prev + change;
        return Math.max(50, Math.min(80, newValue)); // 50-80 dB range
      });
    }, 10000); // Update every 10 seconds

    return () => clearInterval(interval);
  }, []);

  // Load stats for the chosen location
  useEffect(() => {
    loadStats();
  }, [chosenLocation]);

  // Load locations on mount
  useEffect(() => {
    loadLocations();
  }, []);

  // Save threshold to localStorage when it changes
   useEffect(() => {
    localStorage.setItem('noiseThreshold', threshold);
  }, [threshold]);

  const loadLocations = async () => {
    try {
      setLoading(true);
      const response = await locationAPI.getAll();
      const locationData = response.data.locations || [];
      setLocations(locationData);
      
      // Find the chosen location
      const chosen = locationData.find(loc => loc.chosen);
      setChosenLocation(chosen);
    } catch (error) {
      console.error('Error loading locations:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadStats = async () => {
    if (!chosenLocation) return;

    try {
      // Try to get real data from API
      const response = await dataAPI.getAll();
      const allData = response.data || [];
      
      // Filter data for the chosen location
      const roomData = allData.filter(d => d.room_name === chosenLocation.name);
      
      if (roomData.length > 0) {
        // Calculate daily peak and weekly average
        const levels = roomData.map(d => d.sound_level);
        const peak = Math.max(...levels);
        const average = Math.round(levels.reduce((a, b) => a + b, 0) / levels.length);
        
        setDailyPeak(prev => Math.max(prev, peak));
        setWeeklyAverage(average);
      } 
    } 
    catch (error) 
    {
      console.error('Error loading stats:', error);
      // // Fallback to simulated data
    }
  };

  const handleLocationChange = async (locationId) => {
    try {
      await locationAPI.setChosen(locationId);
      await loadLocations();
    } catch (error) {
      console.error('Error changing location:', error);
      alert('Failed to change location');
    }
  };

  const handleThresholdChange = (newThreshold) => {
    setThreshold(newThreshold);
  };

  if (loading) {
    return (
      <div className="min-vh-100 d-flex align-items-center justify-content-center">
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="app-container">
      {/* Header */}
      <header className="app-header text-center mb-4">
        <h1 className="app-title">Kindergarten Noise Meter</h1>
        <p className="app-subtitle text-muted">
          Monitor and analyze classroom noise levels
        </p>
      </header>

      {/* Main Content Grid */}
      <div className="container-fluid">
        <div className="row g-4">
          {/* Left Column - Settings */}
          <div className="col-lg-4">
            <SettingsPanel
              locations={locations}
              chosenLocation={chosenLocation}
              onLocationChange={handleLocationChange}
              threshold={threshold}
              onThresholdChange={handleThresholdChange}
            />
          </div>

        {/* Right Column - Noise Meter */}
        <div className="col-lg-8">
          <NoiseMeter
            currentLevel={currentNoiseLevel}
            threshold={threshold}
            roomName={chosenLocation?.name || 'No Room Selected'}
          />
          </div>
        </div>

      {/* Statistics Cards */}
        <div className="mt-4">
          <StatsCards
          dailyPeak={dailyPeak}
          weeklyAverage={weeklyAverage}
          monitoringRoom={chosenLocation?.name || 'None'}
          // isActive={!!chosenLocation}
          // Though this attribute is in the prototype design, I didn't find it useful for now, in application.
        />
        </div>

        <NoiseAnalytics roomName={chosenLocation?.name || 'None'} 
        allLocations={locations}
        />
      </div>
    </div>
        
  );
}

export default App;