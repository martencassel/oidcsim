package oauth2

type NonceGenerator struct {
}

func (ng *NonceGenerator) Generate() (string, error) {
	return "", nil
}
