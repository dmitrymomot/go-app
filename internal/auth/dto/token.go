package dto

// Token is a data transfer object for token.
// It is used for internal communication between packages.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
