// Package checkgrp provides status check handlers
package checkgrp

import (
	"context"
	"net/http"
	"os"

	db "github.com/qiushiyan/simplebank/business/db/core"
	"github.com/qiushiyan/simplebank/foundation/web"
)

type Handler struct {
	store db.Store
	build string
}

func New(store db.Store, build string) *Handler {
	return &Handler{store: store, build: build}
}

// Readiness godoc
//
//	@Summary		Check the service is ready
//	@Description	Readiness checks if the database is ready and if not will return a 500 status. Do not respond by just returning an error because further up in the call stack it will interpret that as a non-trusted error.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Failure		500
//
//	@Router			/readiness [get]
func (h *Handler) Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	err := h.store.Check(ctx)

	if err != nil {
		return web.NewError(err, http.StatusInternalServerError)
	}

	return web.RespondJsonPlain(ctx, w, "OK", http.StatusOK)
}

type LivenessResponse struct {
	Status     string `json:"status,omitempty"`
	Build      string `json:"build,omitempty"`
	Host       string `json:"host,omitempty"`
	Name       string `json:"name,omitempty"`
	PodIP      string `json:"podIP,omitempty"`
	Node       string `json:"node,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	GOMAXPROCS string `json:"GOMAXPROCS,omitempty"`
}

// Liveness godoc
//
//	@Summary		Check the service is alive
//	@Description	Liveness returns simple status info if the service is alive. If the app is deployed to a Kubernetes cluster, it will also return pod, node, and namespace details via the Downward API. The Kubernetes environment variables need to be set within your Pod/Deployment manifest.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	LivenessResponse
//	@Failure		500
//
//	@Router			/liveness [get]
func (h *Handler) Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := LivenessResponse{
		Status:     "up",
		Build:      h.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: os.Getenv("GOMAXPROCS"),
	}

	// This handler provides a free timer loop.

	return web.RespondJsonPlain(ctx, w, data, http.StatusOK)

}

func (h Handler) Register(a *web.App) {
	a.GET("/readiness", h.Readiness)
	a.GET("/liveness", h.Liveness)
}
