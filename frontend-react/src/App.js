import React, { useState, useEffect, use } from 'react';
import { dataAPI, locationAPI } from './services/api';
import LocationManager from './components/LocationManager';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

function App() {

  useEffect(() => {
    // sample data with getDailySummary method tested and logged to console.
    // Console will show an error if there is no data with that specific room name and time in your database.
    getDailySummary('room 1', new Date(2025, 10, 8)); // nov 7, 2025
  }, []);

  const getDailySummary = async (room, date) => {
    try {
      date.setHours(0, 0, 0, 0);
      date.toISOString();
      const response = await dataAPI.getDailySummary(room, date);
      console.log('Daily Summary:', response.data);
    } catch (error) {
      console.error('Error fetching daily summary:', error);
    }
  };

  return (
    <div className="App">
      <LocationManager />

    </div>
    // Add in here other components later, like routers, navigation bars, etc.
  );
}

export default App;