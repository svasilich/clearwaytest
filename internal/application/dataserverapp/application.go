package dataserverapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	dauth "github.com/svasilich/clearwaytest/domain/auth"
	"github.com/svasilich/clearwaytest/internal/application/asset"
	"github.com/svasilich/clearwaytest/internal/application/auth"
	"github.com/svasilich/clearwaytest/internal/lib/responsehelper"
	"github.com/svasilich/clearwaytest/internal/repository/cwrepo"
)

// DataServerApp is servant for handle user requests.
type DataServerApp struct {
	authorizer    auth.DBUserAuthorizer
	userRetriever auth.DBUserRetriever
	passHasher    auth.Hasher
	assetWriter   asset.DBAssetWriter
}

// NewDataServerApp create instance for DataServer.
func NewDataServerApp(
	authorizer auth.DBUserAuthorizer,
	sessionRetiever auth.DBUserRetriever,
	passHasher auth.Hasher,
	assetWeiter asset.DBAssetWriter,
) *DataServerApp {
	return &DataServerApp{
		authorizer:    authorizer,
		userRetriever: sessionRetiever,
		passHasher:    passHasher,
		assetWriter:   assetWeiter,
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
	if req.Method != http.MethodPost {
		responsehelper.SetupJSONResponse(w, http.StatusBadRequest, "error", "invalid request method")
		return
	}

	// Получить ассет.
	// asset := path.Base(req.RequestURI)
	// log.Default().Println(asset)
	asset, err := getAsset(req.RequestURI)
	if err != nil {
		responsehelper.SetupJSONResponse(w, http.StatusBadRequest, "error", fmt.Errorf("invalid path: %s", req.RequestURI).Error())
		return
	}

	// Проверить, существует ли сессия.
	token, err := getBearerToken(strings.Trim(req.Header.Get("Authorization"), " "))
	if err != nil {
		responsehelper.SetupJSONResponse(w, http.StatusBadRequest, "error", fmt.Errorf("invalid credentials").Error())
		return
	}

	// Получить пользователя.
	uid, err := d.userRetriever.GetUserBySession(req.Context(), dauth.Token(token))
	if err != nil {
		if errors.Is(err, cwrepo.ErrNoOpenSessions) {
			responsehelper.SetupJSONResponse(w, http.StatusUnauthorized, "error", "Session not open or has expired")
			return
		}

		responsehelper.Setup5xx(w, err)
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		responsehelper.Setup5xx(w, err)
		return
	}
	if len(data) == 0 {
		responsehelper.Setup5xx(w, fmt.Errorf("request body is empty"))
		return
	}

	err = d.assetWriter.WriteAsset(req.Context(), asset, uid, data)
	if err != nil {
		responsehelper.Setup5xx(w, fmt.Errorf("can't write data to db: %w", err))
		return
	}

	responsehelper.SetupJSONResponse(w, http.StatusOK, "status", "ok")
}

// Download is handler downloading users data.
func (d *DataServerApp) Download(w http.ResponseWriter, req *http.Request) {
	panic("not implemented")
}
