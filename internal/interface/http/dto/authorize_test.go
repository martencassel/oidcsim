package dto

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseAuthorizeRequest(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    AuthorizeRequest
		wantErr bool
	}{
		{
			name:  "valid request",
			query: "response_type=code&client_id=client1&redirect_uri=http://localhost:8080/callback&scope=read&state=xyz",
			want: AuthorizeRequest{
				ResponseType: "code",
				ClientID:     "client1",
				RedirectURI:  "http://localhost:8080/callback",
				Scope:        "read",
				State:        "xyz",
			},
			wantErr: false,
		},
		{
			name:  "missing client_id",
			query: "response_type=code&redirect_uri=http://localhost:8080/callback&scope=read&state=xyz",
			want: AuthorizeRequest{
				ResponseType: "code",
				RedirectURI:  "http://localhost:8080/callback",
				Scope:        "read",
				State:        "xyz",
			},
			wantErr: false, // Note: Bind does not return error for missing fields
		},
		{
			name:  "missing response_type",
			query: "client_id=client1&redirect_uri=http://localhost:8080/callback&scope=read&state=xyz",
			want: AuthorizeRequest{
				ClientID:    "client1",
				RedirectURI: "http://localhost:8080/callback",
				Scope:       "read",
				State:       "xyz",
			},
			wantErr: false, // Note: Bind does not return error for missing fields
		},
		{
			name:    "empty query",
			query:   "",
			want:    AuthorizeRequest{},
			wantErr: false, // Note: Bind does not return error for empty query
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request, _ = http.NewRequest("GET", "/authorize?"+tt.query, nil)

			var ar AuthorizeRequest
			err := ar.Bind(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ar != tt.want {
				t.Errorf("Bind() got = %+v, want %+v", ar, tt.want)
			}
		})
	}
}
