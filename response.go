package goperiscope

type OAuthRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type OAuthRefreshResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type CreateBroadcastRequest struct {
	Region string `json:"region"`
	Is360  bool   `json:"is_360"`
}

type CreateBroadcastResponse struct {
	Broadcast   Broadcast   `json:"broadcast"`
	VideoAccess VideoAccess `json:"video_access"`
	ShareURL    string      `json:"share_url"`
	Encoder     Encoder     `json:"encoder"`
}

type PublishBroadcastRequest struct {
	BroadcaastID   string `json:"broadcast_id"`
	Title          string `json:"title"`
	ShouldNotTweet bool   `json:"should_not_tweet"`
	Locale         string `json:"locale"`
}

type PublishBroadcastResponse struct {
	Broadcast Broadcast `json:"broadcast"`
}

type StopBroadcastRequest struct {
	BroadcaastID string `json:"broadcast_id"`
}

type DeleteBroadcastRequest struct {
	BroadcaastID string `json:"broadcast_id"`
}

type DeleteBroadcastResponse struct {
}

type GetRegionResponse struct {
	Region string `json:"region"`
}
