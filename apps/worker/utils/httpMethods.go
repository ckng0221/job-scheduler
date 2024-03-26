package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func PatchRequest(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("ADMIN_API_KEY"))
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func PostRequest(url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("ADMIN_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func GetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-api-key", os.Getenv("ADMIN_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
