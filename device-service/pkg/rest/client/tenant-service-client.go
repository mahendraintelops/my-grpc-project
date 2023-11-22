package client

import (
	"io"
	"net/http"
)

func ExecuteNode50() ([]byte, error) {
	response, err := http.Get("tenant-service:8080/ping")

	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}
