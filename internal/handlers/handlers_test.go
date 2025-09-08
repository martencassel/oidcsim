package handlers

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// type AuthorizeRequest struct{}

// // mock service
// type mockAuthorizeService struct {
// 	redirectURL string
// 	err         error
// }

// func (m *mockAuthorizeService) HandleAuthorize(ctx context.Context, req AuthorizeRequest) (string, error) {
// 	return m.redirectURL, m.err
// }

// func TestAuthorizeHandler_ServeHTTP_Success(t *testing.T) {
// 	mockSvc := &mockAuthorizeService{
// 		redirectURL: "https://client.example/cb?code=abc123&state=xyz",
// 	}
// 	h := &AuthorizeHandler{AuthorizeService: nil}
// 	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=cid&redirect_uri=https://client.example/cb&response_type=code&scope=openid&state=xyz", nil)
// 	w := httptest.NewRecorder()
// 	h.ServeHTTP(w, req)

// 	res := w.Result()
// 	defer res.Body.Close()
// 	assert.Equal(t, http.StatusFound, res.StatusCode)
// 	assert.Equal(t, mockSvc.redirectURL, res.Header.Get("Location"))
// }
