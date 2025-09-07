package oauth2

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// logrus
	log "github.com/sirupsen/logrus"

	authapp "github.com/martencassel/oidcsim/internal/application/authentication"
	delegationapp "github.com/martencassel/oidcsim/internal/application/delegation"
	oauth2app "github.com/martencassel/oidcsim/internal/application/oauth2"
	"github.com/martencassel/oidcsim/internal/application/session"
	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	"github.com/martencassel/oidcsim/internal/interface/http/dto"
	middleware "github.com/martencassel/oidcsim/internal/interface/http/middleware"
)

/*

1. Parse -> map to domain request
2. Check authentication context
3. If not logged in -> save request + redirect to /login
4. If logged in -> check consent / delegation
5. Branch on consent decision
6. If consent OK -> run the authorization flow and redirect to client

*/


/*

SessionManager: Abstracts cookie + server‑side state.
TemplateRenderer: Abstracts HTML rendering (could be Go templates, React SSR, etc.).
Domain types: Handlers work with domain value objects, not DTOs from persistence.


*/

type Handler struct {
	sessions      session.SessionManager // interface for session read/write
	AuthSvc       authapp.AuthService
	AuthorizeSvc  oauth2app.AuthorizeService
	DelegationSvc delegationapp.DelegationService
}

func (h *Handler) Authorize(g *gin.Context) {
	ctx := g.Request.Context()
	// 1. Parse query parameters into a DTO
	var dtoReq dto.AuthorizeRequest
	if err := dtoReq.Bind(g); err != nil {
		http.Error(g.Writer, "invalid request", http.StatusBadRequest)
		return
	}
	scopes := fromScopeString(dtoReq.Scope)
	// 2. Map DTO -> domain.AuthorizeRequest (enforce basic invariants)
	domReq := oauth2.AuthorizeRequest{
		ResponseType:        dtoReq.ResponseType,
		ClientID:            dtoReq.ClientID,
		RedirectURI:         dtoReq.RedirectURI,
		Scope:               scopes,
		State:               dtoReq.State,
		CodeChallenge:       dtoReq.CodeChallenge,
		CodeChallengeMethod: dtoReq.CodeChallengeMethod,
		Nonce:               dtoReq.Nonce,
	}
	log.Infof("domReq: %v", domReq)

	// from cookie/session middleware
	sid, _ := middleware.SessionIDFromContext(g.Request.Context())

	// Check current authentication context
	authCtx, ok, _ := h.AuthSvc.Current(g.Request.Context(), sid)

	if !ok || !authCtx.IsValidFor(domReq) {
		_ = h.sessions.SaveAuthorizeRequest(sid, dtoReq)
		// Redirect to login flow
		g.Redirect(http.StatusFound, "/login")
		return
	}
	// 5. Already authenticated → check consent/delegation
	consentResult, err := h.DelegationSvc.EnsureConsent(ctx, authCtx.SubjectID, domReq.ClientID, domReq.Scope)
	if err != nil {
		http.Error(g.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	d := delegationapp.ConsentStatus(consentResult)
	switch d {
	case delegationapp.ConsentRequired:
		// Save request and redirect to consent UI
		_ = h.sessions.SaveAuthorizeRequest(sid, dtoReq)
		g.Redirect(http.StatusFound, "/consent")
		return
	case delegationapp.ConsentDenied:
		http.Error(g.Writer, "consent denied", http.StatusForbidden)
		return
	case delegationapp.ConsentGranted:
		// Proceed to generate authorization code
	default:
		http.Error(g.Writer, "internal error", http.StatusInternalServerError)
		return
	}
	user := oauth2.User{ID: authCtx.SubjectID}
	redirectURL, err := h.AuthorizeSvc.HandleAuthorize(ctx, domReq, user)
	if err != nil {
		http.Error(g.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	// 6. Redirect back to client
	g.Redirect(http.StatusFound, redirectURL)
}

func dtoToDomainAuthorizeRequest(g *gin.Context) dto.AuthorizeRequest {
	var req dto.AuthorizeRequest
	if err := req.Bind(g); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return dto.AuthorizeRequest{}
	}
	return req
}

func fromScopeString(s string) []string {
	if s == "" {
		return []string{}
	}
	return SplitScope(s)
}

func SplitScope(s string) []string {
	var scopes []string
	current := ""
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if current != "" {
				scopes = append(scopes, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		scopes = append(scopes, current)
	}
	return scopes
}

// sessionID := h.sessions.GetID(r)
// authCtx, _ := h.authSvc.Current(ctx, sessionID)
// if !authCtx.IsValidFor(req) {
// 	h.sessions.SaveAuthorizeRequest(sessionID, req)
// 	http.Redirect(w, r, "/login", http.StatusFound)
// 	return
// }
// redirectURL, err := h.authorizeSvc.HandleAuthorize(ctx, req, authCtx.User)
// if err != nil {
// 	h.templates.RenderError(w, err)
// 	return
// }
// http.Redirect(w, r, redirectURL, http.StatusFound)

// func (h *Handler) LoginStart(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	clientID := h.sessions.GetPendingClientID(r)
// 	spec, err := h.authSvc.Initiate(ctx, clientID)
// 	if err != nil {
// 		h.templates.RenderError(w, err)
// 		return
// 	}
// 	step := spec.Steps[0]
// 	ui, _ := h.authSvc.StartStep(ctx, h.sessions.GetID(r), step.Method)
// 	h.templates.RenderStep(w, ui)
// }

// func (h *Handler) LoginStep(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	sessionID := h.sessions.GetID(r)
// 	stepMethod := r.FormValue("method")

// 	inputs := map[string]string{}
// 	_ = r.ParseForm()
// 	for k := range r.PostForm {
// 		inputs[k] = r.PostForm.Get(k)
// 	}

// 	done, err := h.authSvc.CompleteStep(ctx, sessionID, domAuth.AuthMethod(stepMethod), inputs)
// 	if err != nil {
// 		h.templates.RenderError(w, err)
// 		return
// 	}

// 	if !done {
// 		// Render next step
// 		nextStep := h.authSvc.NextStep(sessionID)
// 		ui, _ := h.authSvc.StartStep(ctx, sessionID, nextStep.Method)
// 		h.templates.RenderStep(w, ui)
// 		return
// 	}

// 	// Flow complete → resume /authorize
// 	req := h.sessions.GetAuthorizeRequest(sessionID)
// 	http.Redirect(w, r, "/authorize?"+req.Encode(), http.StatusFound)
// }
