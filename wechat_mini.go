package cbl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	wxAPIPrefix = "https://api.weixin.qq.com"
)

type WXMiniSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

// Code2Session wxlogin code to session id
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func Code2Session(appID string, secret string, code string) (*WXMiniSession, error) {
	api := fmt.Sprintf("%s/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wxAPIPrefix, appID, secret, code)
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, err
	}
	if data.ErrCode != 0 {
		return nil, errors.New(data.ErrMsg)
	}

	return &WXMiniSession{
		OpenID:     data.OpenID,
		SessionKey: data.SessionKey,
		UnionID:    data.UnionID,
	}, nil
}
