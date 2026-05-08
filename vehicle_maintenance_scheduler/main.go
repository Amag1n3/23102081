package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	authURL     = "http://4.224.186.213/evaluation-service/auth"
	logEndpoint = "http://4.224.186.213/evaluation-service/vehicles"
)

var creds = map[string]string{
	"email":        "amoghtyagi22092005@gmail.com",
	"name":         "amogh tyagi",
	"rollNo":       "23102081",
	"accessCode":   "MdprhE",
	"clientID":     "90e96a35-ea8a-4823-8125-4b2a9170a57d",
	"clientSecret": "UybTNQAChyzwYuxB",
}

func getToken() (string, error) {
	b, _ := json.Marshal(creds)
	resp, err := http.Post(authURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("auth parse failed: %w", err)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("empty token received")
	}
	return result.AccessToken, nil
}
