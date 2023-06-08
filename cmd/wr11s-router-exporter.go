package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fcastello/wr11s-router-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type LoginState struct {
	Result int `json:"result,string"`
}

func login(base_url string, username string, password string) error {

	// Check if the username and password are empty
	if username == "" || password == "" {
		return fmt.Errorf("username or password not provided,not all data will be available to the exporter without loging in")
	}

	// Encode the username and password as base64
	encodedUsername := base64.URLEncoding.EncodeToString([]byte(username))
	encodedPassword := base64.URLEncoding.EncodeToString([]byte(password))

	// Create the form data with the encoded username and password
	formData := url.Values{}
	formData.Set("isTest", "false")
	formData.Set("goformId", "LOGIN")
	formData.Set("username", encodedUsername)
	formData.Set("password", encodedPassword)
	loginFormUrl := fmt.Sprintf("%s/goform/goform_set_cmd_process", base_url)
	// Make a POST request to the login endpoint
	resp, err := http.PostForm(loginFormUrl, formData)
	if err != nil {
		return fmt.Errorf("error making login request: %w", err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var loginState LoginState
	err = json.Unmarshal(body, &loginState)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}
	if loginState.Result == 0 {
		return nil
	} else {
		return fmt.Errorf("error loging in with result %v, not all data will be available to the exporter without loging in", loginState.Result)
	}
}

func main() {
	// Initialize logger
	logger := log.New()

	// Set log level (optional)
	logger.SetLevel(log.InfoLevel)
	//Set Logging Options

	logger.SetFormatter(&log.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)
	// Load configuration
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	err := login(cfg.Address, cfg.Username, cfg.Password)
	if err != nil {
		log.Println("Login failed:", err)
	} else {
		log.Println("Login succeded")
	}
}
