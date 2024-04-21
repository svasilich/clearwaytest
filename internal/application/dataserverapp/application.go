package dataserverapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/svasilich/clearwaytest/internal/application/auth"
	"github.com/svasilich/clearwaytest/internal/lib/responsehelper"
)

// DataServerApp is servant for handle user requests.
type DataServerApp struct {
	authorizer       auth.DBUserAuthorizer
	sessionRetriever auth.DBUserSessionRetriever
	passHasher       auth.Hasher
}

// NewDataServerApp create instance for DataServer.
func NewDataServerApp(
	authorizer auth.DBUserAuthorizer,
	sessionRetiever auth.DBUserSessionRetriever,
	passHasher auth.Hasher,
) *DataServerApp {
	return &DataServerApp{
		authorizer:       authorizer,
		sessionRetriever: sessionRetiever,
		passHasher:       passHasher,
	}
}

// Auth is handler for user login.
func (d *DataServerApp) Auth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		responsehelper.SetupJSONResponse(w, http.StatusBadRequest, "error", "invalid request method")
		return
	}

	// Проверить что тело запроса корректно.
	var request authRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	if err != nil {
		responsehelper.SetupJSONResponse(w, http.StatusBadRequest, "error", fmt.Errorf("invalid request: %w", err).Error())
		return
	
	}

	//TODO: проверить логин на валидность.

	passHash, err := d.passHasher(request.Password)
	if err != nil {
		responsehelper.Setup5xx(w, err)
		return
	}

	session, err := d.authorizer.Login(req.Context(), request.Login, passHash)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			responsehelper.SetupJSONResponse(w, http.StatusUnauthorized, "error", "invalid login/password")
			return
		}

		responsehelper.Setup5xx(w, err)
		return
	}

	// TODO: проверить, не истекла ли сессия.

	responsehelper.SetupJSONResponse(w, http.StatusOK, "token", string(session.Token))
}

// Upload is handler uploading users data.
func (d *DataServerApp) Upload(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}

// Download is handler downloading users data.
func (d *DataServerApp) Download(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
