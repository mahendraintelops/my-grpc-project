package client

import (
	"io"
	"net/http"
)

func ExecuteNodeC3() ([]byte, error) {
	response, err := http.Get("user-service:3400/ping")

	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}
