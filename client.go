package goperiscope

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type PeriscopeBuilder struct {
	urlBase      string
	useragent    string
	clientID     string
	clientSecret string
	refreshToken string
}

func NewBuilder(urlBase, useragent, clientID, clientSecret string) PeriscopeBuilder {
	return PeriscopeBuilder{
		urlBase:      urlBase,
		useragent:    useragent,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (b *PeriscopeBuilder) RefreshToken(t string) *PeriscopeBuilder {
	b.refreshToken = t
	return b
}

func (b *PeriscopeBuilder) BuildClient() (Client, error) {

	httpCli := http.Client{
		Timeout: 10 * time.Second,
	}

	authCli := newAuthClient(b.urlBase, &httpCli, b.useragent, b.clientID, b.clientSecret)
	auth, err := authCli.OAuthRefresh(b.refreshToken)
	if err != nil {
		return nil, errors.Wrapf(err, "OAuthRefresh is failed")
	}

	return NewClient(b.urlBase, &httpCli, b.useragent, auth.AccessToken), nil
}

type AuthClient interface {
	OAuthRefresh(refreshToken string) (*OAuthRefreshResponse, error)
}

type AuthClientImpl struct {
	urlBase      string
	httpCli      *http.Client
	useragent    string
	clientID     string
	clientSecret string
}

func newAuthClient(urlBase string, httpCli *http.Client, useragent string, clientID string, clientSecret string) AuthClient {
	return &AuthClientImpl{
		urlBase:      urlBase,
		httpCli:      httpCli,
		useragent:    useragent,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (i AuthClientImpl) OAuthRefresh(refreshToken string) (*OAuthRefreshResponse, error) {
	req := OAuthRefreshRequest{
		GrantType:    "refresh_token",
		ClientID:     i.clientID,
		ClientSecret: i.clientSecret,
		RefreshToken: refreshToken,
	}

	var result OAuthRefreshResponse
	if err := i.request("POST", "/oauth/token", req, &result); err != nil {
		return nil, errors.Wrapf(err, "Periscope /oauth/token is failed")
	}

	return &result, nil
}

func (c AuthClientImpl) request(method, path string, params interface{}, result interface{}) error {

	headers := map[string]string{
		"User-Agent": c.useragent,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	apiURL := fmt.Sprintf("%s%s", c.urlBase, path)
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	// request
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	// error handling for status code
	if resp.StatusCode >= 300 {
		log.Printf("unexpected API Response. statusCode=%d, url=%s", resp.StatusCode, apiURL)

		internalErr := internalError{}

		if err := json.NewDecoder(resp.Body).Decode(&internalErr); err != nil {
			return fmt.Errorf(
				"JSON parse error [statusCode='%d', err='%v']", resp.StatusCode, err,
			)
		}
		return NewError(resp.StatusCode, params.(fmt.Stringer), internalErr)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf(
			"JSON parse error [statusCode='%d', err='%v']", resp.StatusCode, err,
		)
	}

	return nil
}

type Client interface {
	GetRegion() (*GetRegionResponse, error)
	CreateBroadcast(region string, is360 bool) (*CreateBroadcastResponse, error)
	PublishBroadcast(broadcastID string, title string, withTweet bool, locale string) (*PublishBroadcastResponse, error)
	StopBroadcast(broadcastID string) error
	GetBroadcast(broadcastID string) (*Broadcast, error)
	DeleteBroadcast(broadcastID string) error
}

type ClientImpl struct {
	urlBase     string
	httpCli     *http.Client
	useragent   string
	accessToken string
}

func NewClient(urlBase string, httpCli *http.Client, useragent string, accessToken string) Client {
	return &ClientImpl{
		urlBase:     urlBase,
		httpCli:     httpCli,
		useragent:   useragent,
		accessToken: accessToken,
	}
}

func (i ClientImpl) GetRegion() (*GetRegionResponse, error) {

	var result GetRegionResponse
	if err := i.request("GET", "/region", nil, &result); err != nil {
		return nil, errors.Wrapf(err, "Periscope /region is failed")
	}
	return &result, nil
}

func (i ClientImpl) CreateBroadcast(region string, is360 bool) (*CreateBroadcastResponse, error) {

	req := CreateBroadcastRequest{
		Region: region,
		Is360:  is360,
	}

	var result CreateBroadcastResponse
	if err := i.request("POST", "/broadcast/create", req, &result); err != nil {
		return nil, errors.Wrapf(err, "Periscope /broadcast/create is failed")
	}

	return &result, nil
}

func (i ClientImpl) PublishBroadcast(broadcastID string, title string, withTweet bool, locale string) (*PublishBroadcastResponse, error) {

	req := PublishBroadcastRequest{
		BroadcastID:    broadcastID,
		Title:          title,
		ShouldNotTweet: !withTweet,
		Locale:         locale,
	}

	var result PublishBroadcastResponse
	if err := i.request("POST", "/broadcast/publish", req, &result); err != nil {
		return nil, errors.Wrapf(err, "Periscope /broadcast/publish is failed")
	}

	return &result, nil
}

func (i ClientImpl) StopBroadcast(broadcastID string) error {

	req := StopBroadcastRequest{
		BroadcastID: broadcastID,
	}

	if err := i.request("POST", "/broadcast/stop", req, nil); err != nil {
		return errors.Wrapf(err, "Periscope /broadcast/stop is failed")
	}

	return nil
}

func (i ClientImpl) GetBroadcast(broadcastID string) (*Broadcast, error) {
	var result Broadcast
	if err := i.request("GET", fmt.Sprintf("/broadcast?id=%s", broadcastID), nil, &result); err != nil {
		return nil, errors.Wrapf(err, "Periscope /region is failed")
	}
	return &result, nil
}

func (i ClientImpl) DeleteBroadcast(broadcastID string) error {
	req := DeleteBroadcastRequest{
		BroadcastID: broadcastID,
	}

	if err := i.request("POST", "/broadcast/delete", req, nil); err != nil {
		return errors.Wrapf(err, "Periscope /broadcast/delete is failed")
	}

	return nil
}

func (c ClientImpl) request(method, path string, params interface{}, result interface{}) error {

	headers := map[string]string{
		"User-Agent":    c.useragent,
		"Authorization": fmt.Sprintf("Bearer %s", c.accessToken),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	apiURL := fmt.Sprintf("%s%s", c.urlBase, path)
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	// request
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	// error handling for status code
	if resp.StatusCode >= 300 {
		log.Printf("unexpected API Response. statusCode=%d, url=%s", resp.StatusCode, apiURL)

		internalErr := internalError{}

		if err := json.NewDecoder(resp.Body).Decode(&internalErr); err != nil {
			return fmt.Errorf(
				"JSON parse error [statusCode='%d', err='%v']", resp.StatusCode, err,
			)
		}
		return NewError(resp.StatusCode, params.(fmt.Stringer), internalErr)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf(
			"JSON parse error [statusCode='%d', err='%v']", resp.StatusCode, err,
		)
	}

	return nil
}
