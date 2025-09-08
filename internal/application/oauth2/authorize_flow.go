package oauth2

type AuthorizeFlow2 struct {
	Validator AuthorizeValidator
	Handler   AuthorizeHandler
}
