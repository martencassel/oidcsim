package oauth2app

type AuthorizeFlow2 struct {
	Validator AuthorizeValidator
	Handler   AuthorizeHandler
}
