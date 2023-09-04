package app

import (
	"encoding/json"
	"net/http"

	"github.com/NunoFrRibeiro/go_rest_auth/dto"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
	Service "github.com/NunoFrRibeiro/go_rest_auth/service"
)

type AuthHandler struct {
	service Service.AuthService
}

func (h AuthHandler) HandlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "Handler not implemented...")
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("error decodin login request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(loginRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, &token)
		}
	}
}

func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := h.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w, appErr.Code, notAuthrorizedResponse(appErr.Message))
		} else {
			writeResponse(w, http.StatusOK, authrorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthrorizedResponse("missing token"))
	}
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(refreshRequest); err != nil {
		logger.Error("error decoding token request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Refresh(refreshRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, &token)
		}
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func authrorizedResponse() map[string]bool {
	return map[string]bool{
		"isAuthorized": true,
	}
}

func notAuthrorizedResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      message,
	}
}
