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

const API_URL = 'http://localhost:8080/api';
// Adjust the URL as needed
const AUTH = btoa('kids_noisemeter_admin:passwordkids');
const locations = [];

//Load locations on page load
window.addEventListener('DOMContentLoaded', loadLocations);

// fetch locations from backend
function loadLocations() {
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
            locations.length = 0; // Clear locations array
            data.locations.forEach(location => locations.push(location));
            loadLocationsList();

            // --> function for old dropdown menu
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
            // <-- function for old dropdown menu
        })
        .catch(error => {
            console.error('Error fetching locations:', error);
            alert('Failed to load locations. Make sure the backend server is running.');
            // this alert will also show if the locations table is empty
        });
}

// create a dropdown menu with delete buttons for locations 
function loadLocationsList() {
    const list = document.getElementById("locationList");
    const toggle = document.getElementById("toggleBtn");
    list.innerHTML = ''; // Clear list
    if (locations.length > 0) toggle.innerHTML = ''; // Clear button text

    locations.forEach(location => {
        const li = document.createElement("li");
        li.setAttribute("role", "option");
        li.dataset.id = location.id;
        li.className = "location-item list-element";
        li.innerHTML = `
        <span class="locationName text-center">${location.name}</span>
        <button class="delete" title="Delete ${location.name}">X</button>
        `;

        if (location.chosen) {
            toggle.textContent = location.name + " ▾";
            console.log("Chosen location loaded: " + location.name);
        }

        li.addEventListener("click", () => {
            toggle.textContent = location.name + " ▾";
            changeChosenLocation(location.id);
            list.hidden = true;
        });
        li.querySelector(".delete").addEventListener("click", (event) => {
            event.stopPropagation();
            deleteLocation(location.id, location.name);
        });
        list.appendChild(li);
    });
}

// handle deleting a location
function deleteLocation(id, name) {
    if (!confirm(`Are you sure you want to delete location (${name})?`)) return;

    // checking location exists
    const index = locations.findIndex(location => location.id === id);
    if (index === -1) {
        console.error(`Location with id ${id} not found`);
        return;
    }

    // removes deleted locations from the array
    const [removed] = locations.splice(index, 1);
    loadLocationsList();

    // calls the deletion request in the backend
    fetch(`${API_URL}/locations/${id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Basic ${AUTH}`,
            'Content-Type': 'application/json'
        }
    })
    .then(response => { if (!response.ok) throw new Error('Failed to delete location'); })
    .catch(error => {
        locations.splice(index, 0, removed);
        loadLocationsList();
        alert('Failed to delete location; location restored');
        console.error(error);
    }); 
}

// opening dropdown menu
document.getElementById("toggleBtn").addEventListener("click", () => {
    const list = document.getElementById("locationList");
    list.hidden = !list.hidden;
});

// Handle adding a new location
document.getElementById("addLocationForm").addEventListener("submit", function (event) {
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
//document.getElementById("toggleBtn").addEventListener("change", function (event) {
    //const locationId = event.target.value;
function changeChosenLocation(locationId) {

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
};
