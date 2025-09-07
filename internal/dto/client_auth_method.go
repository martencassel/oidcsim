package dto

type ClientAuthMethod string

const (
	ClientSecretBasic ClientAuthMethod = "client_secret_basic"
	ClientSecretPost  ClientAuthMethod = "client_secret_post"
	PrivateKeyJWT     ClientAuthMethod = "private_key_jwt"
	TLSClientAuth     ClientAuthMethod = "tls_client_auth"
	None              ClientAuthMethod = "none"
)
