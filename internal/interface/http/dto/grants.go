package dto

type GrantListResponse struct {
	Grants []Grant `json:"grants"`
}

type Grant struct {
	GrantID   string   `json:"grant_id"`
	ClientID  string   `json:"client_id"`
	Scope     []string `json:"scope"`
	CreatedAt int64    `json:"created_at"`
}
