package health

import (
	"gopher-social-backend-server/pkg/utils"
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) GetRouterHealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
	}

	if err := utils.WriteJSON(w, http.StatusOK, response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}
