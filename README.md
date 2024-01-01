# Go Ajax-WASM Example

## Description

Go Ajax-WASM Example is a cutting-edge project demonstrating the power of WebAssembly (WASM) in conjunction with Golang to simplify AJAX calls in JavaScript. By leveraging WASM modules written in Go, this repository showcases how to perform server-side operations seamlessly within a web application, bypassing the complexities of traditional AJAX in JS.

## Installation

### Prerequisites
- Ensure you have [sqlc](https://sqlc.dev/) installed for database interaction.
- Install [mkcert](https://github.com/FiloSottile/mkcert) for local HTTPS development.
- [TinyGo](https://tinygo.org/) is required to build the WASM module.
- [Golang](https://go.dev/) should be installed to build and run the server.

## Steps
1. Generate Database Interactions

```
sqlc generate
```

2. Create a Local Certificate for HTTPS

```
mkcert -cert-file cmd/webserver/server-cert.pem -key-file cmd/webserver/server-key.pem localhost 127.0.0.1 ::1
```

3. Build the WASM Module

```
tinygo build -o static/authors.wasm -target=wasm cmd/wasm/main.go
```

4. Build the Server

```
go build -o server.exe .\cmd\webserver\main.go
```

5. Run the Server

```
server.exe
```

After running the server, the website will be accessible at https://localhost:3000.

## Usage

Interacting with WebAssembly Functions
Example JavaScript function to interact with the WASM module:

```
// Fetches all authors from the server's database using a WASM function.
function getAllAuthors() {
    return new Promise((resolve, reject) => {
        wasmGetAllAuthors((result, err) => {
            if (err) reject(err);
            else resolve(result);
        });
    });
}
```

## Contributing
Contributions to the Go Ajax-WASM Example are welcome. Please ensure to follow standard coding practices and provide tests for new features.

## License


## Contact
For support or to report issues, please create an issue on the GitHub repository.