

/**
 * Interacts with a WebAssembly function to fetch all authors from the server's internal database.
 * This function utilizes the 'wasmGetAllAuthors' function defined in the WebAssembly module.
 * @returns A Promise that resolves with all authors from the database or rejects with an error.
 */
function getAllAuthors() {
    // Return a new Promise object that will handle the asynchronous WebAssembly call.
    return new Promise((resolve, reject) => {
        // Call 'wasmGetAllAuthors', a WebAssembly function to fetch all authors.
        wasmGetAllAuthors((result, err) => {
            // If an error occurs during the WebAssembly call, reject the Promise with the error.
            if (err) reject(err);
            // If the call is successful, resolve the Promise with the fetched authors.
            else resolve(result);
        });
    });
}

/**
 * Creates a new author in the server's internal database through a WebAssembly call.
 * This function takes two parameters 'name' and 'bio' and utilizes the 'wasmCreateAuthor' function from the WebAssembly module.
 * @param name: the name of the author.
 * @param bio: the biography of the author.
 * @returns A Promise that resolves with the newly created author or rejects with an error.
 */
function createAuthor(name, bio) {
    // Return a new Promise object to handle the asynchronous creation of an author.
    return new Promise((resolve, reject) => {
        // Call 'wasmCreateAuthor', a WebAssembly function to create a new author.
        wasmCreateAuthor(name, bio, (result, err) => {
            // If an error occurs during the WebAssembly call, reject the Promise with the error.
            if (err) reject(err);
            // If the call is successful, resolve the Promise with the newly created author.
            else resolve(result);
        });
    });
}
