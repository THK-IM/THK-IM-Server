package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	msgApiPostUserOnlineStatusUrl = "/user/online"
)

type (
	MsgApi interface {
		PostUserOnlineStatus(req dto.PostUserOnlineReq) error
	}

	defaultMsgApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *http.Client
	}
)

func (d defaultMsgApi) PostUserOnlineStatus(req dto.PostUserOnlineReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, msgApiPostUserOnlineStatusUrl)
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

func NewMsgApi(sdk conf.Sdk, logger *logrus.Entry) MsgApi {
	return defaultMsgApi{
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
