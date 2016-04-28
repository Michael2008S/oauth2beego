// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package apps

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"

	"github.com/Michael2008S/oauth2beego"
)

type QQ struct {
	BaseProvider
}

func (p *QQ) GetType() social.SocialType {
	return social.SocialQQ
}

func (p *QQ) GetName() string {
	return "QQ"
}

func (p *QQ) GetPath() string {
	return "qq"
}

func (p *QQ) GetIndentify(tok *social.Token) (string, error) {
	uri := "https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	noCallBackStr := bodyStr[10 : len(bodyStr)-4]

	var f interface{}
	err = json.Unmarshal([]byte(noCallBackStr), &f)
	if err != nil {
		beego.Error(err)

		vals, err := url.ParseQuery(body)
		if err != nil {
			return "", err
		}
		fmt.Println("noCallBackStr:", noCallBackStr)
		if vals.Get("code") != "" {
			return "", fmt.Errorf("code: %s, msg: %s", vals.Get("code"), vals.Get("msg"))
		}

	} else {
		m := f.(map[string]interface{})
		beego.Info(m)
		beego.Error("openid:", (m["openid"]).(string))
		return (m["openid"]).(string), nil
	}
	return "", fmt.Errorf("openid no found")

}

var _ social.Provider = new(QQ)

func NewQQ(clientId, secret string) *QQ {
	p := new(QQ)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://graph.qq.com/oauth2.0/authorize"
	p.TokenURL = "https://graph.qq.com/oauth2.0/token"
	p.RedirectURL = social.DefaultAppUrl + "loginsocial/qq/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
