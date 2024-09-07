package types

type Session struct {
	Username        string `json:"username"`
	IsAdmin         bool   `json:"is_admin"`
	Timestamp       int64  `json:"timestamp"`
	IsAuthenticated bool   `json:"is_authenticated"`
}
