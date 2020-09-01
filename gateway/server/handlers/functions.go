package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HandleFunctionCall creates a functions request endpoint
func HandleFunctionCall(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the path parameters
		vars := mux.Vars(r)
		projectID := vars["project"]
		service := vars["service"]
		function := vars["func"]

		auth, err := modules.Auth(projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		functions, err := modules.Functions(projectID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Load the params from the body
		req := model.FunctionsRequest{}
		_ = json.NewDecoder(r.Body).Decode(&req)
		defer utils.CloseTheCloser(r.Body)

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Set a default timeout value
		if req.Timeout == 0 {
			req.Timeout = 60 // set default context to 10 second
		}

		// Create a new context
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(req.Timeout)*time.Second)
		defer cancel()

		actions, reqParams, err := auth.IsFuncCallAuthorised(ctx, projectID, service, function, token, req.Params)
		if err != nil {
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err.Error())
			return
		}

		status, result, err := functions.CallWithContext(ctx, service, function, token, reqParams, req.Params)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Receieved error from service call (%s:%s)", service, function), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err.Error())
			return
		}

		_ = auth.PostProcessMethod(ctx, actions, result)

		_ = helpers.Response.SendResponse(ctx, w, status, result)
	}
}
