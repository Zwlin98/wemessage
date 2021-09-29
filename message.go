package wemessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Response interface {
	UnMarshalFromJSON([]byte) error
}

type Message interface {
	ToJSON() ([]byte, error)
}

type MessageResponse struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	MsgID        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}

func (r *MessageResponse) UnMarshalFromJSON(p []byte) error {
	return json.Unmarshal(p, r)
}

func GetSendURL(c *Client) (string, error) {
	accessToken, err := c.AccessToken()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/cgi-bin/message/send?access_token=%s", c.BaseURL, accessToken)
	if err != nil {
		return "", err
	}

	return url, nil
}

func SendMessage(c *Client, m Message) (*MessageResponse, error) {
	url, err := GetSendURL(c)
	if err != nil {
		return nil, err
	}

	body, err := m.ToJSON()
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(url, "application/json", bytes.NewReader(body))

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r MessageResponse
	err = r.UnMarshalFromJSON(body)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

type TextMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"` //Must be text
	AgentID int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe                   int `json:"safe"`
	EnableIdTrans          int `json:"enable_id_trans"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m TextMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type ImageMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Image   struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
	Safe                   int `json:"safe"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m ImageMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type VoiceMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Voice   struct {
		MediaID string `json:"media_id"`
	} `json:"voice"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m VoiceMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type VideoMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Video   struct {
		MediaID     string `json:"media_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"video"`
	Safe                   int `json:"safe"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m VideoMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type FileMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	File    struct {
		MediaID string `json:"media_id"`
	} `json:"file"`
	Safe                   int `json:"safe"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m FileMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type TextCardMessage struct {
	ToUser   string `json:"touser"`
	ToParty  string `json:"toparty"`
	ToTag    string `json:"totag"`
	MsgType  string `json:"msgtype"`
	AgentID  int    `json:"agentid"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		BtnText     string `json:"btntxt"`
	} `json:"textcard"`
	EnableIDTrans          int `json:"enable_id_trans"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m TextCardMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type NewsMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	News    struct {
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			PicURL      string `json:"picurl"`
			Appid       string `json:"appid"`
			PagePath    string `json:"pagepath"`
		} `json:"articles"`
	} `json:"news"`
	EnableIDTrans          int `json:"enable_id_trans"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m NewsMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

type MarkdownMessage struct {
	ToUser   string `json:"touser"`
	ToParty  string `json:"toparty"`
	ToTag    string `json:"totag"`
	MsgType  string `json:"msgtype"`
	AgentID  int    `json:"agentid"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
	EnableDuplicateCheck   int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

func (m MarkdownMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
