package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"aspera-terminal-ui/internal/config"
)

type DownloadCommandsResponse struct {
	ResultCode int `json:"resultCode"`
	Results    struct {
		AscpArgs     string `json:"ascpargs"`
		CompleteDate string `json:"CompleteDate"`
	} `json:"results"`
	Message string `json:"message"`
}

type FileItem struct {
	FileName string `json:"FileName"`
	Size     int64  `json:"Size"`
}

type RequestDetail struct {
	RequestId  string     `json:"RequestId"`
	Subject    string     `json:"Subject"`
	Files      []FileItem `json:"Files"`
	Mail       MailDetail `json:"Mail"`
	Status     string     `json:"Status"`
	ExpireTime string     `json:"ExpireTime"`
}

type RequestDetailResponse struct {
	ResultCode int           `json:"resultCode"`
	Results    RequestDetail `json:"results"`
	Message    string        `json:"message"`
}

// GetDownloadCommand queries the Aspera API for the download command.
// If filenames is not empty, it adds them to the "filenames" header separated by /.
func GetDownloadCommand(requestID string, filenames []string) (*DownloadCommandsResponse, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	token, err := EnsureValidToken()
	if err != nil {
		return nil, err
	}

	encodedReqID := url.PathEscape(requestID)
	apiURL := fmt.Sprintf("%s/Customer/FEX/P2P/v1/requests/%s/AsperaDownloadCommands", cfg.Endpoint, encodedReqID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	
	if len(filenames) > 0 {
		req.Header.Set("filenames", url.PathEscape(strings.Join(filenames, "/")))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("token expired or unauthorized (401)")
	}
	
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch download commands, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var data DownloadCommandsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if data.Results.AscpArgs == "" {
		return nil, fmt.Errorf("ascpargs missing from response: %s", data.Message)
	}

	return &data, nil
}

// GetRequestDetail returns the details of a specific request including its file list.
func GetRequestDetail(requestID string) (*RequestDetail, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	token, err := EnsureValidToken()
	if err != nil {
		return nil, err
	}

	encodedReqID := url.PathEscape(requestID)
	apiURL := fmt.Sprintf("%s/Customer/FEX/P2P/v1/requests/%s", cfg.Endpoint, encodedReqID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch request detail, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Printf("DEBUG: Request Detail Raw: %s\n", string(bodyBytes))

	var data RequestDetailResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v, raw: %s", err, string(bodyBytes))
	}

	return &data.Results, nil
}
