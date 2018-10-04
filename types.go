package goperiscope

import "fmt"

type Broadcast struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Title string `json:"title"`
}

type Encoder struct {
	StreamKey                string              `json:"stream_key"`
	RtmpURL                  string              `json:"rtmp_url"`
	RtmpsURL                 string              `json:"rtmps_url"`
	DisplayName              string              `json:"display_name"`
	RecommendedConfiguration StreamConfiguration `json:"recommended_configuration"`
	IsStreamActive           bool                `json:"is_stream_active"`
}

type StreamConfiguration struct {
	VideoCodec        string `json:"video_codec"`
	VideoBitrate      uint32 `json:"video_bitrate"`
	Framerate         uint32 `json:"framerate"`
	KeyframeInterval  uint32 `json:"keyframe_interval"`
	Width             uint32 `json:"width"`
	Height            uint32 `json:"height"`
	AudioCodec        string `json:"audio_codec"`
	AudioSamplingRate uint32 `json:"audio_sampling_rate"`
	AudioBitrate      uint32 `json:"audio_bitrate"`
	AudioNumChannels  uint32 `json:"audio_num_channels"`
}

type User struct {
	ID               string             `json:"id"`
	Username         string             `json:"username"`
	TwitterID        string             `json:"twitter_id"`
	TwitterUsername  string             `json:"twitter_username"`
	Description      string             `json:"description"`
	DisplayName      string             `json:"display_name"`
	ProfileImageURLs []ProfileImageURLs `json:"profile_image_urls"`
}

type ProfileImageURLs struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
	SslURL string `json:"ssl_url"`
	URL    string `json:"url"`
}

type VideoAccess struct {
	HlsURL      string `json:"hls_url"`
	HTTPSHlsURL string `json:"https_hls_url"`
}

type ChatMessage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Text  string `json:"text"`
	User  User   `json:"user"`
	Color string `json:"color"`
}

type HeartMessage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	User  User   `json:"user"`
	Color string `json:"color"`
}

type JoinMessage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	User  User   `json:"user"`
	Color string `json:"color"`
}

type ScreenshotMessage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	User  User   `json:"user"`
	Color string `json:"color"`
}

type ShareMessage struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Service string `json:"service"`
	User    User   `json:"user"`
	Color   string `json:"color"`
}

type SuperHeartMessage struct {
	Type   string `json:"type"`
	User   User   `json:"user"`
	Color  string `json:"color"`
	Amount int32  `json:"amount"`
	Tier   int32  `json:"tier"`
}

type ViewerCountMessage struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Live  int32  `json:"live"`
	Total int32  `json:"total"`
}

type ErrorMessage struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type internalError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (i internalError) String() string {
	return fmt.Sprintf("message=%v, documentationURL=%v", i.Message, i.DocumentationURL)
}

type Error struct {
	StatusCode    int
	Params        fmt.Stringer
	InternalError fmt.Stringer
}

func (e Error) Error() string {
	return fmt.Sprintf(
		`statusCode="%d" params="%s" error="%v"]`,
		e.StatusCode,
		e.Params,
		e.InternalError,
	)
}

func (e Error) HTTPStatusCode() int {
	return e.StatusCode
}

func NewError(statusCode int, params, internalErr fmt.Stringer) error {
	return &Error{
		StatusCode:    statusCode,
		Params:        params,
		InternalError: internalErr,
	}
}
