package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ListFilter struct {
	PageIndex int `json:"PageIndex"`
	PageSize  int `json:"PageSize"`
}

type SharedRequest struct {
	RequestId  string `json:"RequestId"`
	Subject    string `json:"Subject"`
	ShareTime  string `json:"ShareTime"`
	ExpireTime string `json:"ExpireTime"`
	Status     string `json:"Status"`
	System     string `json:"System"`
}

type ReceivedRequest struct {
	MailId     string `json:"MailId"`
	UserEmail  string `json:"UserEmail"`
	System     string `json:"System"`
	SendTime   string `json:"SendTime"`
	ExpireTime string `json:"ExpireTime"`
	Title      string `json:"Title"`
	Status     string `json:"Status"`
	RequestId  string `json:"RequestId"`
	Downloaded bool   `json:"Downloaded"`
}

type SharedResponse struct {
	Count    int             `json:"Count"`
	DataList []SharedRequest `json:"DataList"`
}

type ReceivedResponse struct {
	Count    int               `json:"Count"`
	DataList []ReceivedRequest `json:"DataList"`
}

func QuerySharedRequests(endpoint, token string, filter ListFilter) (*SharedResponse, error) {
	validToken, err := EnsureValidToken()
	if err == nil {
		token = validToken
	}

	reqURL := fmt.Sprintf("%s/Customer/FEX/P2P/v1/sharedRequestsQuerier", strings.TrimRight(endpoint, "/"))

	b, _ := json.Marshal(filter)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var wrapper struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Results SharedResponse  `json:"results"`
	}

	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.Results, nil
}

func QueryReceivedRequests(endpoint, token string, filter ListFilter) (*ReceivedResponse, error) {
	validToken, err := EnsureValidToken()
	if err == nil {
		token = validToken
	}

	reqURL := fmt.Sprintf("%s/Customer/FEX/P2P/v1/receivedRequestsQuerier", strings.TrimRight(endpoint, "/"))

	b, _ := json.Marshal(filter)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var wrapper struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Results ReceivedResponse `json:"results"`
	}

	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return nil, err
	}

	return &wrapper.Results, nil
}
