package authentication

import (
	"gopher-social-backend-server/pkg/mailer"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
	mailer.SendActivationEmail(user.Email, token)

	if err := utils.WriteJSON(w, http.StatusCreated, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthenticationHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, "token is required")
		return
	}

	email, err := utils.VerifyActivationToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.AuthenticationStore.GetUserByEmail(email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if user.IsActivated {
		utils.WriteError(w, http.StatusBadRequest, "user is already activated")
		return
	}

	user.IsActivated = true
	if err := h.AuthenticationStore.UpdateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mailer.SendAccountActivatedEmail(user.Email)

	if err := utils.WriteJSON(w, http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthenticationHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload userLoginPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.AuthenticationStore.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if !user.IsActivated {
		utils.WriteError(w, http.StatusUnauthorized, "account is not activated")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	accessToken, expirationTime := utils.GenerateAccessToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    accessToken,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"token": accessToken}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthenticationHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthenticationHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var payload userForgotPasswordPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.AuthenticationStore.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	token := utils.GeneratePasswordResetToken(user.Email)
	mailer.SendPasswordResetEmail(user.Email, token)

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "password reset email sent"}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthenticationHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, "token is required")
		return
	}

	email, err := utils.VerifyPasswordResetToken(token)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload userPasswordResetPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.AuthenticationStore.GetUserByEmail(email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), BCRYPT_COST)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = string(hashedPassword)
	if err := h.AuthenticationStore.UpdateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	mailer.SendPasswordChangedEmail(user.Email)

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "password changed"}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
