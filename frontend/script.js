/*
// For creating a locations table in the backend database
CREATE TABLE IF NOT EXISTS locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    chosen BOOLEAN NOT NULL DEFAULT 0 CHECK (chosen IN (0, 1))
);

// for ensuring no more than one location is chosen at a time
CREATE UNIQUE INDEX IF NOT EXISTS only_one_chosen_location ON locations(id) WHERE chosen = 1;

// for changing the chosen location
BEGIN IMMEDIATE;
UPDATE locations SET chosen = 0 WHERE chosen = 1;
UPDATE locations SET chosen = 1 WHERE id = ?;
COMMIT;

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

// handle adding a new location
document.getElementById("addLocationForm").addEventListener("submit", function (event) {
    event.preventDefault();

    const location = document.getElementById("newLocation").value;
    console.log(`Location to add: ${location}`);

    // send location to db
});

// handle changing the chosen location
document.getElementById("locationSelect").addEventListener("change", function (event) {
    const locationId = event.target.value;
    console.log(`Location ID to be chosen: ${locationId}`);

    // Change chosen location based on locationId
});