package test

import (
	"net/http"
	"time"

	"ccu/api"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
)

// GetTest godoc
// @Summary      Gets a test value from the service, sanity check
// @Description  Will ask the service to generate a test json and return it back to the requester
// @Tags         Tests
// @Accept       json
// @Produce      json
// @Success      200  {array}   mAPI.Test
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /test-no-auth [get]
func GetTest(w http.ResponseWriter, r *http.Request) {
	log.Info("In test handler -------------------------")
	data := mAPI.Test{
		Id:          "TEST",
		DateCreated: time.Now(),
		Amount:      1,
		Usd:         2,
		Change:      3.0,
	}

	api.RespondOK(w, data)
}
