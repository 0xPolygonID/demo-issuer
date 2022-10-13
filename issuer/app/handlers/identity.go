package handlers

import (
	"github.com/go-chi/render"
	"lightissuer/services"
	"net/http"
)

type IdentityHandler struct {
	identityService *services.IdentityService
}

func NewIdentityHandler(service *services.IdentityService) *IdentityHandler {
	return &IdentityHandler{service}
}

// GetIdentity POST /api/v1/identity
func (i *IdentityHandler) GetIdentity(w http.ResponseWriter, r *http.Request) {
	identity := i.identityService.Identity
	render.JSON(w, r, identity)
}
