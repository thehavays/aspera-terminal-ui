package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"aspera-terminal-ui/internal/config"
)

type ShareRequest struct {
	FileHashAlgorithm string   `json:"FileHashAlgorithm"`
	Files             []FileInfo `json:"Files"`
	Mail              MailInfo `json:"Mail"`
	Subject           string   `json:"Subject"`
	System            string   `json:"System"`
}

type FileInfo struct {
	FileName string `json:"FileName"`
	Size     int64  `json:"Size"`
	Hash     string `json:"Hash,omitempty"`
}

type Recipient struct {
	Email string `json:"Email"`
}

type MailInfo struct {
	BccList []Recipient `json:"BccList"`
	Body    string      `json:"Body"`
	CcList  []Recipient `json:"CcList"`
	Title   string      `json:"Title"`
	ToList  []Recipient `json:"ToList"`
}

type MailDetail struct {
	BccList []string `json:"BccList"`
	Body    string   `json:"Body"`
	CcList  []string `json:"CcList"`
	Title   string   `json:"Title"`
	ToList  []string `json:"ToList"`
}

type ShareResponse struct {
	ResultCode int `json:"resultCode"`
	Results    struct {
		AscpArgs  string `json:"ascpargs"`
		RequestId string `json:"RequestId"`
	} `json:"results"`
	Message string `json:"message"`
}

// CreateShareRequest creates a share request and returns ascp arguments and request ID.
func CreateShareRequest(subject, body string, toList, ccList, bccList []string, files []string) (*ShareResponse, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	token, err := EnsureValidToken()
	if err != nil {
		return nil, err
	}

	reqBody := ShareRequest{
		FileHashAlgorithm: "", // We can add SHA1 later if needed
		Mail: MailInfo{
			BccList: toRecipients(bccList),
			Body:    body,
			CcList:  toRecipients(ccList),
			Title:   subject,
			ToList:  toRecipients(toList),
		},
		Subject: subject,
		System:  "API",
	}

	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %v", f, err)
		}
		reqBody.Files = append(reqBody.Files, FileInfo{
			FileName: filepath.Base(f),
			Size:     info.Size(),
		})
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/Customer/FEX/P2P/v1/sharedRequests", cfg.Endpoint)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBytes))
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

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to create share request, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var data ShareResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if data.ResultCode != 0 {
		return nil, fmt.Errorf("API error: %s", data.Message)
	}

	// Check if the inner results object has a non-zero code (business logic error)
	var rawResults struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
	}
	// Try to unmarshal just the results part to check for errors
	var rawMap map[string]json.RawMessage
	json.Unmarshal(bodyBytes, &rawMap)
	if resultsData, ok := rawMap["results"]; ok {
		if err := json.Unmarshal(resultsData, &rawResults); err == nil {
			if rawResults.Code != 0 && data.Results.RequestId == "" {
				return nil, fmt.Errorf("server constraint: %s", rawResults.Content)
			}
		}
	}

	if data.Results.RequestId == "" {
		return nil, fmt.Errorf("RequestId missing in results. Response: %s", string(bodyBytes))
	}

	return &data, nil
}

func toRecipients(emails []string) []Recipient {
	res := make([]Recipient, 0, len(emails))
	for _, e := range emails {
		res = append(res, Recipient{Email: e})
	}
	return res
}
