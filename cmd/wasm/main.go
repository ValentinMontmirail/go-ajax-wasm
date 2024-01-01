package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"syscall/js"
	"time"
	"wasm/pkg/authors"
	"wasm/pkg/token"
)

func WasmRequest(method string, url string, body []byte, statusCodeExpected int) ([]byte, error) {
	fmt.Println("[" + method + "] - " + url + " - IN")

	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Errorf("Could not create the request: %s", err.Error())
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Requested-By", token.REQUESTED_BY) // Custom header for identification

	// Create a custom HTTP client
	client := &http.Client{Timeout: time.Second * 30}

	res, err := client.Do(r)
	if err != nil {
		fmt.Errorf("Could not perform the request: %s", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != statusCodeExpected {
		fmt.Errorf("Expected %v | Got %v", statusCodeExpected, res.StatusCode)
		return nil, err
	}

	output, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Could not read the body: %s", err.Error())
		return nil, err
	}

	fmt.Println("[" + method + "] - " + url + " - OUT")
	return output, nil
}

// Convert any struct to a js.Value object dynamically
func structToJSValue(s interface{}) js.Value {
	val := reflect.ValueOf(s)
	typ := val.Type()

	// Create a new JavaScript object
	jsObj := js.Global().Get("Object").New()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name := typ.Field(i).Name

		// Convert the Go value to a JavaScript value
		var jsVal js.Value
		switch field.Kind() {
		case reflect.Bool:
			jsVal = js.ValueOf(field.Bool())
		case reflect.String:
			jsVal = js.ValueOf(field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			jsVal = js.ValueOf(field.Int())
		case reflect.Float32, reflect.Float64:
			jsVal = js.ValueOf(field.Float())
		default:
			// Skip or handle unsupported types
			continue
		}

		// Set the field in the JavaScript object
		jsObj.Set(name, jsVal)
	}

	return jsObj
}

func GetAllAuthors(this js.Value, args []js.Value) interface{} {

	if len(args) == 0 {
		err := errors.New("Function must be called from JavaScript")
		fmt.Errorf("%s", err.Error())
		return err
	}

	// Extract the callback function from the arguments
	// This is assuming that the function is being called from JavaScript, and a callback is passed as the first argument
	callback := args[0]

	err := fetchToken()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		callback.Invoke(js.Null(), fmt.Sprintf("Token creation error: %v", err))
		return err
	}

	// Launch a new Goroutine. This allows the HTTP request to be handled asynchronously, so it doesn't block the main thread.
	go func() {
		body, err := WasmRequest("GET", "/api/v1/authors", nil, http.StatusOK)
		if err != nil {
			fmt.Errorf("%s", err.Error())
			callback.Invoke(js.Null(), err.Error())
			return
		}

		// Parse the JSON body into a slice of Author structs
		var authors []authors.Author
		err = json.Unmarshal(body, &authors)
		if err != nil {
			callback.Invoke(js.Null(), err.Error())
			return
		}

		// Convert the slice of authors to a JavaScript array
		jsAuthors := js.Global().Get("Array").New(len(authors))
		for i, author := range authors {
			jsAuthors.SetIndex(i, structToJSValue(author))
		}

		callback.Invoke(jsAuthors, js.Null())
	}()

	// Since the HTTP request is handled in a separate Goroutine, return immediately
	// The actual response will be handled via the callback
	return nil
}

// Fetch token from server with authentication
func fetchToken() error {
	_, err := WasmRequest("GET", "/api/v1/token", nil, http.StatusOK)
	if err != nil {
		fmt.Errorf("%s", err.Error())
	}
	return err
}

func CreateAuthor(this js.Value, args []js.Value) interface{} {

	// Extract name and bio from arguments
	name := args[0].String()
	bio := args[1].String()
	callback := args[2]

	if len(args) != 3 {
		err := errors.New("Function must be called from JavaScript with 2 string parameters")
		fmt.Errorf("%s", err.Error())
		return err
	}

	// // Construct the JSON body with the provided name and bio
	jsonBody, err := json.Marshal(authors.CreateAuthorParams{
		Name: name,
		Bio:  bio,
	})

	if err != nil {
		fmt.Errorf("%s", err.Error())
		callback.Invoke(js.Null(), err.Error())
		return nil
	}

	err = fetchToken()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		callback.Invoke(js.Null(), err.Error())
		return err
	}

	go func() {
		body, err := WasmRequest("POST", "/api/v1/authors", jsonBody, http.StatusOK)
		if err != nil {
			fmt.Errorf("%s", err.Error())
			callback.Invoke(js.Null(), err.Error())
			return
		}
		callback.Invoke(string(body), js.Null())
	}()

	return nil
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("wasmCreateAuthor", js.FuncOf(CreateAuthor))
	js.Global().Set("wasmGetAllAuthors", js.FuncOf(GetAllAuthors))
	<-done
}
