package witai

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const witAPIURL = "https://api.wit.ai/message"

// Client is used to interact with WitAI.
type Client struct {
	OAuthToken string
	Version    string
}

// NewClient returns a newly initialized witAI client with the given OAuth token.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("Invalid Wit token given")
	}

	return &Client{
		Version:    "16/09/2017",
		OAuthToken: token,
	}, nil
}

func (c *Client) newHTTPClient() *http.Client {
	tlscfg := &tls.Config{}

	tlscfg.InsecureSkipVerify = true

	transport := &http.Transport{
		TLSClientConfig: tlscfg,
	}

	return &http.Client{Transport: transport}
}

func (c *Client) newRequest(message string) (*http.Request, error) {
	params := url.Values{}
	params.Set("q", message)
	params.Set("v", c.Version)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s?%s", witAPIURL, params.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.OAuthToken))

	return req, nil
}

func (c *Client) decodeMeaning(responseBody io.ReadCloser) (*Meaning, error) {
	defer responseBody.Close()

	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))

	newMeaning := &Meaning{}
	err = json.Unmarshal(body, newMeaning)

	return newMeaning, err
}

// GetMeaning returns the estimated intent for the given message.
func (c *Client) GetMeaning(message string) (*Meaning, error) {
	httpClient := c.newHTTPClient()
	req, err := c.newRequest(message)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return c.decodeMeaning(res.Body)
}
