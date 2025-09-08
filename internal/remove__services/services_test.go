package services

// func TestHandleAuthorize_BuildsRedirect(t *testing.T) {
// 	registry := authflows.NewRegistry()
// 	svc := NewAuthorizeServiceImpl(registry)
// 	req := dto.AuthorizeRequest{
// 		ClientID:     "cid",
// 		RedirectURI:  "https://client.example/cb",
// 		ResponseType: "code",
// 		Scope:        []string{"openid", "profile", "email"},
// 		State:        "xyz",
// 	}
// 	redirectURL, err := svc.HandleAuthorize(context.Background(), req)
// 	assert.NoError(t, err)
// 	u, _ := url.Parse(redirectURL)
// 	q := u.Query()
// 	assert.Equal(t, "xyz", q.Get("state"))
// 	assert.NotEmpty(t, q.Get("code"))
// 	assert.True(t, strings.HasPrefix(redirectURL, req.RedirectURI))
// }
