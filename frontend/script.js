/*

//I deleted some lines of comments here,
as they seemed to be just memos and hints on how to implement stuff.

// for adding a new location
BEGIN IMMEDIATE;
UPDATE locations SET chosen = 0 WHERE chosen = 1;
INSERT INTO locations (name, chosen) VALUES (?, 1);
COMMIT;
*/

/*
// getting the locations for the dropdown menu (will eventually be modified to fit the backend)
// /api/locations = could be the API endpoint for locations
fetch('/api/locations', {
    method: 'GET',
    // (I think the browser was having error message
    // becasue it lacks headers definition here.)
})
.then(response => response.json())
.then(data => {
    const locationSelect = document.getElementById("locationSelect");
    data.locations.forEach(location => {
        const option = document.createElement("option");
        option.value = location.id;
        option.text = location.name;
        if (location.chosen) {
            option.selected = true;
        }
        locationSelect.add(option);
    });
})
.catch(error => console.error('Error fetching locations:', error));
*/

// // handle adding a new location
// document.getElementById("addLocationForm").addEventListener("submit", function (event) {
//     event.preventDefault();

//     const location = document.getElementById("newLocation").value;
//     console.log(`Location to add: ${location}`);

//     // send location to db
// });

// // handle changing the chosen location
// document.getElementById("locationSelect").addEventListener("change", function (event) {
//     const locationId = event.target.value;
//     console.log(`Location ID to be chosen: ${locationId}`);

//     // Change chosen location based on locationId
// });

const API_URL = 'http://localhost:8080'; 
// Adjust the URL as needed
const AUTH = btoa('kids_noisemeter_admin:passwordkids');

//Load locations on page load
window.addEventListener('DOMContentLoaded', loadLocations);

function loadLocations()
{
    fetch(`${API_URL}/locations`,
        {
            method: 'GET',
            headers: {
                'Authorization': `Basic ${AUTH}`,
                'Content-Type': 'application/json'
                // (I think the browser was having error message
                // becasue it lacks headers definition here.)
            }
        })
        .then(response => {
            if (!response.ok) throw new Error('Failed to fetch locations');
            return response.json();
        })
        .then(data => {
            const locationSelect = document.getElementById("locationSelect");
            locationSelect.innerHTML = ''; // Clear existing options
            
            data.locations.forEach(location => {
                const option = document.createElement("option");
                option.value = location.id;
                option.text = location.name;
                if (location.chosen) {
                    option.selected = true;
                }
                locationSelect.add(option);
            });
        })
        .catch(error => {
            console.error('Error fetching locations:', error);
            alert('Failed to load locations. Make sure the backend server is running.');
        });
}

// Handle adding a new location
document.getElementById("addLocationForm").addEventListener("submit", function (event){
    event.preventDefault();

    const locationName = document.getElementById("newLocation").value;

    fetch(`${API_URL}/locations`, 
    {
        method: 'POST',
        headers:
        {
            'Authorization': `Basic ${AUTH}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify
        ({ 
            name: locationName,
            chosen: true 
        })
    })

    .then(response => {
        if (!response.ok) throw new Error('Failed to add location');
        return response.json();
    })

    .then(data => {
        console.log('Location added: ', data);
        document.getElementById("newLocation").value = '';
        // Clear input
        loadLocations();
        // Reload locations list
    })

    .catch(error => {
        console.error('Error adding location:', error);
        alert('Failed to add location.');
    });
});

// Handle changing the chosen location
document.getElementById("locationSelect").addEventListener("change", function (event) {
    const locationId = event.target.value;
    
    fetch(`${API_URL}/locations/${locationId}`, {
        method: 'PUT',
        headers: {
            'Authorization': `Basic ${AUTH}`,
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) throw new Error('Failed to update location');
        return response.json();
    })
    .then(data => {
        console.log('Location updated:', data);
    })
    .catch(error => {
        console.error('Error updating location:', error);
        alert('Failed to update location.');
    });
});
