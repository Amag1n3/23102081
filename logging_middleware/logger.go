package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	endpoint = "http://4.224.186.213/evaluation-service/logs"
	token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiYXVkIjoiaHR0cDovLzIwLjI0NC41Ni4xNDQvZXZhbHVhdGlvbi1zZXJ2aWNlIiwiZW1haWwiOiJhbW9naHR5YWdpMjIwOTIwMDVAZ21haWwuY29tIiwiZXhwIjoxNzc4MjMzMjk2LCJpYXQiOjE3NzgyMzIzOTYsImlzcyI6IkFmZm9yZCBNZWRpY2FsIFRlY2hub2xvZ2llcyBQcml2YXRlIExpbWl0ZWQiLCJqdGkiOiI0OWNiYmYzYi03ZjI0LTQyMTgtOTM1Mi02MzcxZGU1MzQ2YWQiLCJsb2NhbGUiOiJlbi1JTiIsIm5hbWUiOiJhbW9naCB0eWFnaSIsInN1YiI6IjkwZTk2YTM1LWVhOGEtNDgyMy04MTI1LTRiMmE5MTcwYTU3ZCJ9LCJlbWFpbCI6ImFtb2dodHlhZ2kyMjA5MjAwNUBnbWFpbC5jb20iLCJuYW1lIjoiYW1vZ2ggdHlhZ2kiLCJyb2xsTm8iOiIyMzEwMjA4MSIsImFjY2Vzc0NvZGUiOiJNZHByaEUiLCJjbGllbnRJRCI6IjkwZTk2YTM1LWVhOGEtNDgyMy04MTI1LTRiMmE5MTcwYTU3ZCIsImNsaWVudFNlY3JldCI6IlV5YlROUUFDaHl6d1l1eEIifQ.e-PnajMP9m2Ft1Vh87ZS7FWEBDmCJLG4v7wm5ddr_EI"
)

type req struct {
	Stack   string `json:"stack"`
	Level   string `json:"level"`
	Package string `json:"package"`
	Message string `json:"message"`
}

func Log(stack, level, pack, message string) error {
	b, err := json.Marshal(req{stack, level, pack, message})
	if err != nil {
		return fmt.Errorf("logger: marshal failed: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("logger: couldn't build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("logger: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logger: server error %d: %s", resp.StatusCode, body)
	}
	return nil
}
