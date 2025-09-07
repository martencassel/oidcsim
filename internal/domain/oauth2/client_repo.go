package oauth2

type ClientRepo interface {
	FindByID(id string) (*Client, error)
}
