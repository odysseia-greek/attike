package handlers

import (
	"github.com/odysseia-greek/plato/middleware"
	"github.com/odysseia-greek/plato/models"
	"net/http"
	"time"
)

func Health(w http.ResponseWriter, req *http.Request) {
	healthy := models.Health{
		Time: time.Now().String(),
	}

	if !healthy.Healthy {
		middleware.ResponseWithCustomCode(w, http.StatusBadGateway, healthy)
		return
	}
	middleware.ResponseWithJson(w, healthy)
}
