package authentication

import (
	"gopher-social-backend-server/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

var BCRYPT_COST = utils.GetEnvAsInt("BCRYPT_COST", 10)

type AuthenticationHandler struct {
	AuthenticationStore AuthenticationStore
}

func (h *AuthenticationHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload userRegisterPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), BCRYPT_COST)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  string(hashedPassword),
	}
	if err := h.AuthenticationStore.CreateUser(&user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token := utils.GenerateActivationToken(user.Email)
	utils.SendActivationEmail(user.Email, token)

	if err := utils.WriteJSON(w, http.StatusCreated, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
