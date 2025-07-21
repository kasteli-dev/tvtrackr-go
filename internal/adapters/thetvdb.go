package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kasteli-dev/tvtrackr-go/internal/config"
)

type TheTVDBClient struct {
	APIKey     string
	Token      string
	HTTPClient *http.Client
}

func NewTheTVDBClient(apiKey string) *TheTVDBClient {
	return &TheTVDBClient{
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: config.HTTPTimeout},
	}
}

// Authenticate and get JWT token
func (c *TheTVDBClient) Authenticate() error {
	url := config.TheTVDBLoginURL
	body := map[string]string{"apikey": c.APIKey}
	jsonBody, _ := json.Marshal(body)
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("thetvdb auth failed: %s %s", resp.Status, string(b))
	}

	var respBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}
	c.Token = respBody.Data.Token
	return nil
}

// Search series by name
func (c *TheTVDBClient) SearchSeries(query string) ([]map[string]interface{}, error) {
	if c.Token == "" {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
	}
	url := fmt.Sprintf(config.TheTVDBSearchURL, query)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("thetvdb search failed: %s %s", resp.Status, string(b))
	}
	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
