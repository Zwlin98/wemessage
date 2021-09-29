package wemessage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	CorpID      string
	CorpSecret  string
	accessToken string
	expireTime  time.Time
	BaseURL     string
	HttpClient  *http.Client
}

type TokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *TokenResponse) UnMarshalFromJSON(p []byte) error {
	return json.Unmarshal(p, r)
}

func NewClient(corpID, corpSecret string) *Client {
	c := &Client{
		CorpID:     corpID,
		CorpSecret: corpSecret,
		BaseURL:    "https://qyapi.weixin.qq.com",
		HttpClient: &http.Client{},
	}
	return c
}

func (c Client) accessTokenURL() string {
	return fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.BaseURL, c.CorpID, c.CorpSecret)
}

func (c *Client) AccessToken() (string, error) {
	if c.TokenExpired() {
		err := c.RenewAccessToken()
		if err != nil {
			return "", err
		}
	}
	return c.accessToken, nil
}

func (c *Client) RenewAccessToken() error {
	if !c.TokenExpired() {
		return nil
	}
	url := c.accessTokenURL()
	r, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var resp TokenResponse
	err = resp.UnMarshalFromJSON(body)
	if err != nil {
		return err
	}

	if resp.ErrCode != 0 {
		return errors.New(resp.ErrMsg)
	}

	c.accessToken = resp.AccessToken
	c.expireTime = time.Now().Add(time.Duration(resp.ExpiresIn/2) * time.Second)

	return nil
}

func (c Client) TokenExpired() bool {
	return time.Now().After(c.expireTime)
}

func (c *Client) Get(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}
	return c.HttpClient.Do(req)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, err
	}
	return c.HttpClient.Do(req)
}
