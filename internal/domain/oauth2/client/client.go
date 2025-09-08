package client

// Client entity + redirect URI validation
type Client struct {
	ID           string
	RedirectURIs []string
	Secret       string
}

func (c Client) IsRedirectURIMatching(uri string) bool {
	for _, r := range c.RedirectURIs {
		if r == uri {
			return true
		}
	}
	return false
}

func (c Client) ValidateRedirectURI(uri string) bool {
	return c.IsRedirectURIMatching(uri)
}

func (c Client) AllowsRedirect(uri string) bool {
	return c.ValidateRedirectURI(uri)
}
