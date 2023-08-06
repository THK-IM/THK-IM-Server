package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	userApiGetUserIdByTokenUrl string = "/user"
	userApiOnlineStatusUrl     string = "/user/online"
	contentType                string = "application/json"
)

type (
	PostUserOnlineReq struct {
		UserId   int64 `json:"user_id"`
		IsOnline bool  `json:"is_online"`
	}

	GetUserIdByTokenReq struct {
		Platform string `json:"platform"`
		Token    string `json:"token"`
	}

	GetUserIdByTokenRes struct {
		UserId int64 `json:"user_id"`
	}

	UserApi interface {
		PostUserOnlineStatus(req PostUserOnlineReq) error
		GetUserIdByToken(req GetUserIdByTokenReq) (*GetUserIdByTokenRes, error)
	}

	defaultUserApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *http.Client
	}
)

func (d defaultUserApi) PostUserOnlineStatus(req PostUserOnlineReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, userApiOnlineStatusUrl)
	bodyReader := bytes.NewReader(dataBytes)
	res, e := d.client.Post(url, contentType, bodyReader)
	if e != nil {
		d.logger.Error(e)
		return e
	}
	if res.StatusCode != http.StatusOK {
		body := make([]byte, 0)
		if _, e = res.Body.Read(body); e == nil {
			e = errors.New(string(body))
			d.logger.Error(e)
			return e
		} else {
			e = errors.New(fmt.Sprintf("StatusCode:%d", res.StatusCode))
			d.logger.Error(e)
			return e
		}
	} else {
		return nil
	}
}

func (d defaultUserApi) GetUserIdByToken(req GetUserIdByTokenReq) (*GetUserIdByTokenRes, error) {
	url := fmt.Sprintf("%s%s?token=%s&platform=%s", d.endpoint, userApiGetUserIdByTokenUrl, req.Token, req.Platform)
	res, e := d.client.Get(url)
	if e != nil {
		d.logger.Error(e)
		return nil, e
	}
	if res.StatusCode != http.StatusOK {
		body := make([]byte, 0)
		if _, e = res.Body.Read(body); e == nil {
			e = errors.New(string(body))
			d.logger.Error(e)
			return nil, e
		} else {
			e = errors.New(fmt.Sprintf("StatusCode:%d", res.StatusCode))
			d.logger.Error(e)
			return nil, e
		}
	} else {
		body := make([]byte, 0)
		if _, e = res.Body.Read(body); e == nil {
			tokenRes := &GetUserIdByTokenRes{}
			e = json.Unmarshal(body, tokenRes)
			if e != nil {
				d.logger.Error(e)
				return nil, e
			} else {
				return tokenRes, nil
			}
		} else {
			e = errors.New("BodyReadError")
			d.logger.Error(e)
			return nil, e
		}
	}
}

func NewUserApi(sdk conf.Sdk, logger *logrus.Entry) UserApi {
	return defaultUserApi{
		endpoint: sdk.Endpoint,
		logger:   logger.WithField("rpc", sdk.Name),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        30,
				MaxIdleConnsPerHost: 30,
				IdleConnTimeout:     5 * time.Second,
			},
			Timeout: 5 * time.Second,
		},
	}
}
