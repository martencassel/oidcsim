package http

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) LoginStart(g *gin.Context) {
	_ = g.Request.Context()
	// sessionID, _ := middleware.SessionIDFromContext(ctx)
	// clientID := h.Sessions.GetPendingClientID(sessionID)
	// spec, err := h.AuthSvc.Initiate(ctx, clientID)
	// if err != nil {
	// 	g.String(http.StatusInternalServerError, "auth initiation failed: %v", err)
	// 	return
	// }
	// step := spec.Steps[0]
	// ui, err := h.AuthSvc.StartStep(ctx, sessionID, step.Method)
	// if err != nil {
	// 	g.String(http.StatusInternalServerError, "start step failed: %v", err)
	// 	return
	// }
	// // Render login UI (could be HTML template or JSON)
	// g.HTML(http.StatusOK, "login_step.tmpl", gin.H{
	// 	"method": step.Method,
	// 	"ui":     ui,
	// })
}

// POST
func (h *Handler) LoginStep(g *gin.Context) {
	_ = g.Request.Context()
	// sessionID, _ := middleware.SessionIDFromContext(ctx)
	// stepMethod := g.PostForm("method")
	// inputs := map[string]string{}
	// for k := range g.Request.PostForm {
	// 	inputs[k] = g.PostForm(k)
	// }
	// done, err := h.AuthSvc.CompleteStep(ctx, sessionID, stepMethod, inputs)
	// if err != nil {
	// 	g.String(http.StatusUnauthorized, "authentication failed: %v", err)
	// 	return
	// }
	// if !done {
	// 	nextStep := h.AuthSvc.NextStep(sessionID)
	// 	ui, _ := h.AuthSvc.StartStep(ctx, sessionID, nextStep.Method)
	// 	g.HTML(http.StatusOK, "login_step.tmpl", gin.H{
	// 		"method": nextStep.Method,
	// 		"ui":     ui,
	// 	})
	// 	return
	// }
	// // Flow complete → resume /authorize
	// req := h.Sessions.GetAuthorizeRequest(sessionID)
	//g.Redirect(http.StatusFound, "/authorize?"+req.Encode())
}
