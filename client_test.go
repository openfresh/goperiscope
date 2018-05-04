package goperiscope

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuthRefresh(t *testing.T) {

	refreshToken := "hoge_refresh_token"
	clientID := "hoge_client_id"
	clientSecret := "hoge_client_secret"

	method := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/oauth/token"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		params := OAuthRefreshRequest{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "refresh_token", params.GrantType)
		assert.Equal(t, clientID, params.ClientID)
		assert.Equal(t, clientSecret, params.ClientSecret)
		assert.Equal(t, refreshToken, params.RefreshToken)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "new_token",
			"user": {
				"id": "1111",
				"twitter_username": "hoge111",
				"twitter_id": "11111111",
				"username": "hoge_username",
				"display_name": "hoge_display_name",
				"description": "hoge description",
				"profile_image_urls": [{
					"url": "http://example.com/small.png",
					"ssl_url": "https://example.com/small.png",
					"width": 128,
					"height": 90
				}]
			},
			"expires_in": 15551999,
			"token_type": "Bearer"
		}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := AuthClientImpl{
		urlBase:      ts.URL,
		httpCli:      httpCli,
		useragent:    "goperiscope test",
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	result, err := c.OAuthRefresh(refreshToken)
	assert.NoError(t, err)

	assert.Equal(t, "new_token", result.AccessToken)
	assert.Equal(t, "1111", result.User.ID)
	assert.Equal(t, "hoge111", result.User.TwitterUsername)
	assert.Equal(t, "11111111", result.User.TwitterID)
	assert.Equal(t, "hoge_display_name", result.User.DisplayName)
	assert.Equal(t, "hoge description", result.User.Description)
	assert.Equal(t, "http://example.com/small.png", result.User.ProfileImageURLs[0].URL)
	assert.Equal(t, "https://example.com/small.png", result.User.ProfileImageURLs[0].SslURL)
	assert.Equal(t, uint32(128), result.User.ProfileImageURLs[0].Width)
	assert.Equal(t, uint32(90), result.User.ProfileImageURLs[0].Height)
	assert.Equal(t, 15551999, result.ExpiresIn)
	assert.Equal(t, "Bearer", result.TokenType)
}

func TestCreateBroadcast(t *testing.T) {

	method := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/broadcast/create"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		params := CreateBroadcastRequest{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "ap-northeast-1", params.Region)
		assert.Equal(t, false, params.Is360)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"broadcast": {
				"id": "hogehogehoge",
				"state": "not_started",
				"title": ""
			},
			"video_access": {
				"hls_url": "https://api.pscp.tv/v1/hls?token=hogehogehoge"
			},
			"share_url": "https://www.pscp.tv/w/hogehogehoge",
			"encoder": {
				"stream_key": "hoge_stream_key",
				"display_name": "hoge_display_name",
				"rtmp_url": "hoge_rtmp_url",
				"rtmps_url": "hoge_rtmps_url",
				"recommended_configuration": {
					"video_codec": "H.264/AVC",
					"video_bitrate": 800000,
					"framerate": 30,
					"keyframe_interval": 3,
					"width": 960,
					"height": 540,
					"audio_codec": "AAC",
					"audio_sampling_rate": 44100,
					"audio_bitrate": 96000,
					"audio_num_channels": 2
				},
				"is_stream_active": false
			}
		}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := ClientImpl{
		urlBase:     ts.URL,
		httpCli:     httpCli,
		useragent:   "goperiscope test",
		accessToken: "test-token",
	}

	result, err := c.CreateBroadcast("ap-northeast-1", false)
	assert.NoError(t, err)

	assert.Equal(t, "hogehogehoge", result.Broadcast.ID)
	assert.Equal(t, "not_started", result.Broadcast.State)
	assert.Equal(t, "https://api.pscp.tv/v1/hls?token=hogehogehoge", result.VideoAccess.HlsURL)
	assert.Equal(t, "https://www.pscp.tv/w/hogehogehoge", result.ShareURL)
	assert.Equal(t, "hoge_stream_key", result.Encoder.StreamKey)
	assert.Equal(t, "hoge_display_name", result.Encoder.DisplayName)
	assert.Equal(t, "hoge_rtmp_url", result.Encoder.RtmpURL)
	assert.Equal(t, "hoge_rtmps_url", result.Encoder.RtmpsURL)
	assert.Equal(t, "H.264/AVC", result.Encoder.RecommendedConfiguration.VideoCodec)
	assert.Equal(t, false, result.Encoder.IsStreamActive)
}

func TestGetBroadcast(t *testing.T) {

	method := "GET"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/broadcast"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		assert.Equal(t, "broadcast_id", r.URL.Query().Get("id"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"broadcast_id","state":"not_started","title":"title"}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := ClientImpl{
		urlBase:     ts.URL,
		httpCli:     httpCli,
		useragent:   "goperiscope test",
		accessToken: "test-token",
	}

	result, err := c.GetBroadcast("broadcast_id")
	assert.NoError(t, err)
	assert.Equal(t, "broadcast_id", result.ID)
	assert.Equal(t, "not_started", result.State)
	assert.Equal(t, "title", result.Title)
}

func TestPublishBroadcast(t *testing.T) {
	method := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/broadcast/publish"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		params := PublishBroadcastRequest{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "broadcast_id", params.BroadcaastID)
		assert.Equal(t, "title", params.Title)
		assert.Equal(t, true, params.ShouldNotTweet)
		assert.Equal(t, "ja_JP", params.Locale)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"broadcast":{"id":"broadcast_id","state":"running","title":"title"}}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := ClientImpl{
		urlBase:     ts.URL,
		httpCli:     httpCli,
		useragent:   "goperiscope test",
		accessToken: "test-token",
	}

	result, err := c.PublishBroadcast("broadcast_id", "title", false, "ja_JP")
	assert.NoError(t, err)
	assert.Equal(t, "broadcast_id", result.Broadcast.ID)
	assert.Equal(t, "running", result.Broadcast.State)
	assert.Equal(t, "title", result.Broadcast.Title)
}

func TestStopBroadcast(t *testing.T) {
	method := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/broadcast/stop"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		params := StopBroadcastRequest{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "broadcast_id", params.BroadcaastID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := ClientImpl{
		urlBase:     ts.URL,
		httpCli:     httpCli,
		useragent:   "goperiscope test",
		accessToken: "test-token",
	}

	err := c.StopBroadcast("broadcast_id")
	assert.NoError(t, err)
}

func TestDeleteBroadcast(t *testing.T) {

	method := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("r.Method = '%s', want '%s'", r.Method, method)
		}

		correctPath := "/broadcast/delete"
		if r.URL.Path != correctPath {
			t.Errorf("r.URL.Path ='%v', want '%v'", r.URL.Path, correctPath)
		}

		params := DeleteBroadcastRequest{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "broadcast_id", params.BroadcaastID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	httpCli := &http.Client{}

	c := ClientImpl{
		urlBase:     ts.URL,
		httpCli:     httpCli,
		useragent:   "goperiscope test",
		accessToken: "test-token",
	}

	err := c.DeleteBroadcast("broadcast_id")
	assert.NoError(t, err)
}
