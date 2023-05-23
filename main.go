package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type LoginState struct {
	Result int `json:"result,string"`
}

func login() error {
	username := os.Getenv("WR11_USER")
	password := os.Getenv("WR11_PASSWORD")

	// Check if the username and password are empty
	if username == "" || password == "" {
		return fmt.Errorf("username or password not provided")
	}

	// Encode the username and password as base64
	encodedUsername := base64.URLEncoding.EncodeToString([]byte(username))
	encodedPassword := base64.URLEncoding.EncodeToString([]byte(password))
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println(encodedUsername)
	fmt.Println(encodedPassword)

	// Create the form data with the encoded username and password
	formData := url.Values{}
	formData.Set("isTest", "false")
	formData.Set("goformId", "LOGIN")
	formData.Set("username", encodedUsername)
	formData.Set("password", encodedPassword)

	// Make a POST request to the login endpoint
	resp, err := http.PostForm("http://192.168.150.1/goform/goform_set_cmd_process", formData)
	if err != nil {
		return fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	fmt.Println(string(body))

	var loginState LoginState
	err = json.Unmarshal(body, &loginState)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}
	fmt.Println(loginState.Result)

	return nil
}

func main() {
	err := login()
	if err != nil {
		fmt.Println("Login failed:", err)
	}
}
