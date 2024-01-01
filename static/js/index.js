function updateTable() {
    getAllAuthors().then(response => {
        let data = response;

        const tableBody = document.getElementById('authorsTable').getElementsByTagName('tbody')[0];
        tableBody.innerHTML = ''; // Clear existing rows

        data.forEach(author => {
            const row = tableBody.insertRow();
            const idCell = row.insertCell(0);
            const nameCell = row.insertCell(1);
            const bioCell = row.insertCell(2);

            idCell.textContent = author.ID;
            nameCell.textContent = author.Name;
            bioCell.textContent = author.Bio;
        });
    }).catch(err => {
        showToast("Error: Server refused to update list of authors")
    });
}

// Function to generate random name and bio
function getRandomAuthorData() {
    const names = ["Alice", "Bob", "Charlie", "Diana"];
    const bios = ["Author of fantasy novels", "Renowned historian", "Science fiction writer", "Award-winning journalist"];

    const randomName = names[Math.floor(Math.random() * names.length)];
    const randomBio = bios[Math.floor(Math.random() * bios.length)];

    return { name: randomName, bio: randomBio };
}

// Event listener for the button
document.getElementById("createAuthorBtn").addEventListener("click", function() {
    const { name, bio } = getRandomAuthorData();
    createAuthor(name, bio)
        .then(author => {
            console.log("Author created:", author);
            updateTable();
        })
        .catch(err => showToast("Error creating an author"));
});


function showToast(content = "Unknown error") { //You can change the default value
    // Get the snackbar DIV
    var x = document.getElementById("snackbar");
    
    //Change the text (not mandatory, but I think you might be willing to do it)
    x.innerHTML = content;
  
    // Add the "show" class to DIV
    x.className = "show";
  
    // After 3 seconds, remove the show class from DIV
    setTimeout(function(){ x.className = x.className.replace("show", ""); }, 3000);
}

updateTable();