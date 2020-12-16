package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

type API struct {
	URL string
}

const (
	APP_KEY    = "minchi"
	APP_SECRET = "golang-practice"
)

type AccessToken struct {
	Token string `json:"token"`
}

func NewAPI(url string) *API {
	return &API{URL: url}
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, fmt.Sprintf("%s/%s", a.URL, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (a *API) getAccessToken(ctx context.Context) (string, error) {
	data := url.Values{"app_key": {APP_KEY}, "app_secret": {APP_SECRET}}
	body := strings.NewReader(data.Encode())
	resp, err := http.Post(fmt.Sprintf("%s/%s", a.URL, "auth"), "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", err
	}

	resBody, _ := ioutil.ReadAll(resp.Body)
	var accessToken AccessToken
	_ = json.Unmarshal(resBody, &accessToken)
	return accessToken.Token, nil
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s?name=%s", a.URL, "api/v1/tags", name), nil)
	req.Header.Add("token", token)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
