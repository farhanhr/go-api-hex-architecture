package response

type SuccessfulAuthResponse struct {
	Meta
	AccessToken string `json:"acess_token"`
	ExpiredAt int64 `json:"expired_at"`
}