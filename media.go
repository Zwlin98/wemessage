package wemessage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type TemporaryMediaResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func (r *TemporaryMediaResponse) UnMarshalFromJSON(p []byte) error {
	return json.Unmarshal(p, r)
}

const (
	IMAGE        = "image"
	VOICE        = "voice"
	VIDEO        = "video"
	FILE         = "file"
	MinFileSize  = 5
	MaxImageSize = 2 * 1024 * 1024
	MaxVoiceSize = 2 * 1024 * 1024
	MaxVideoSize = 10 * 1024 * 1024
	MaxFileSize  = 20 * 1024 * 1024
)

var ErrSizeLimit = errors.New("exceed the file size limit")

func GetTemporaryMediaURL(c *Client, mediaType string) (string, error) {
	accessToken, err := c.AccessToken()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/cgi-bin/media/upload?access_token=%s&type=%s", c.BaseURL, accessToken, mediaType), nil
}

func UploadTemporaryMedia(c *Client, mediaType string, name string) (string, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	switch mediaType {
	case IMAGE:
		if fi.Size() > MaxImageSize || fi.Size() < MinFileSize {
			return "", ErrSizeLimit
		}
		suffix := path.Ext(name)
		if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
			return "", errors.New("invalid image format")
		}
	case VOICE:
		if fi.Size() > MaxVoiceSize || fi.Size() < MinFileSize {
			return "", ErrSizeLimit
		}
		suffix := path.Ext(name)
		if suffix != ".amr" {
			return "", errors.New("invalid voice format")
		}
	case VIDEO:
		if fi.Size() > MaxVideoSize || fi.Size() < MinFileSize {
			return "", ErrSizeLimit
		}
		suffix := path.Ext(name)
		if suffix != ".mp4" {
			return "", errors.New("invalid video format")
		}
	case FILE:
		if fi.Size() > MaxFileSize || fi.Size() < MinFileSize {
			return "", ErrSizeLimit
		}
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	fw, err := w.CreateFormFile("media", file.Name())
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return "", err
	}
	_ = w.Close()
	url, err := GetTemporaryMediaURL(c, mediaType)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := c.HttpClient.Do(req)

	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resp TemporaryMediaResponse
	err = json.Unmarshal(bin, &resp)
	if err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", fmt.Errorf("ErrCode:%d, %s", resp.ErrCode, resp.ErrMsg)
	}
	return resp.MediaId, nil
}
