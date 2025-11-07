import React, { useState, useEffect, use } from 'react';
import { locationAPI } from '../services/api';

// was going to start creating a component for building the chart(s) xD
function Charts() {
    const [locations, setLocations] = useState([]);
    useEffect(() => {
        const fetchLocations = async () => {
            try {
                const response = await locationAPI.getAll();
                setLocations(response.data.locations || []);
            } catch (error) {
                console.error('Error fetching locations:', error);
            }
        };
        fetchLocations();
    }, []);


}

export default { Charts };