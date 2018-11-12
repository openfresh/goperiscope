package goperiscope

import "fmt"

type OAuthRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func (r OAuthRefreshRequest) String() string {
	return fmt.Sprintf("grant_type=%s,client_id=%s,client_secret=%s,refresh_token=%s", r.GrantType, r.ClientID, r.ClientSecret, r.RefreshToken)
}

type OAuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (r OAuthRefreshResponse) String() string {
	return fmt.Sprintf("access_token=%s,user=[%s],expires_in=%d,token_type=%s", r.AccessToken, r.User.String(), r.ExpiresIn, r.TokenType)
}

type CreateBroadcastRequest struct {
	Region       string `json:"region"`
	Is360        bool   `json:"is_360"`
	IsLowLatency bool   `json:"is_low_latency"`
}

func (r CreateBroadcastRequest) String() string {
	return fmt.Sprintf("region=%s,is_360=%t,is_low_latency=%t", r.Region, r.Is360, r.IsLowLatency)
}

type CreateBroadcastResponse struct {
	Broadcast   Broadcast   `json:"broadcast"`
	VideoAccess VideoAccess `json:"video_access"`
	ShareURL    string      `json:"share_url"`
	Encoder     Encoder     `json:"encoder"`
}

func (r CreateBroadcastResponse) String() string {
	return fmt.Sprintf("broadcast={%s},video_access={%s},share_url=%s,encoder=%s",
		r.Broadcast.String(), r.VideoAccess.String(), r.ShareURL, r.Encoder.String())
}

type PublishBroadcastRequest struct {
	BroadcastID       string `json:"broadcast_id"`
	Title             string `json:"title"`
	ShouldNotTweet    bool   `json:"should_not_tweet"`
	Locale            string `json:"locale"`
	EnableSuperHearts bool   `json:"enable_super_hearts"`
}

func (r PublishBroadcastRequest) String() string {
	return fmt.Sprintf("broadcast_id=%s,title=%s,should_not_tweet=%t,locale=%s,enable_super_hearts=%t",
		r.BroadcastID, r.Title, r.ShouldNotTweet, r.Locale, r.EnableSuperHearts)
}

type PublishBroadcastResponse struct {
	Broadcast Broadcast `json:"broadcast"`
}

func (r PublishBroadcastResponse) String() string {
	return fmt.Sprintf("broadcast=%s", r.Broadcast.String())
}

type StopBroadcastRequest struct {
	BroadcastID string `json:"broadcast_id"`
}

func (r StopBroadcastRequest) String() string {
	return fmt.Sprintf("broadcast_id=%s", r.BroadcastID)
}

type DeleteBroadcastRequest struct {
	BroadcastID string `json:"broadcast_id"`
}

func (r DeleteBroadcastRequest) String() string {
	return fmt.Sprintf("broadcast_id=%s", r.BroadcastID)
}

type DeleteBroadcastResponse struct {
}

func (r DeleteBroadcastResponse) String() string {
	return ""
}

type GetRegionResponse struct {
	Region string `json:"region"`
}

func (r GetRegionResponse) String() string {
	return fmt.Sprintf("region=%s", r.Region)
}
