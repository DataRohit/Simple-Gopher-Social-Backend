package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/mailer"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var validate = validator.New()

var BCRYPT_COST = utils.GetEnvAsInt("BCRYPT_COST", 10)
var OAUTH_STATE = utils.GetEnvAsString("OAUTH_STATE", "randomstate")

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
		OAuth:     constants.ProviderNone,
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

	if user.OAuth != constants.ProviderNone {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Sprintf("login with %s", user.OAuth))
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

func (h *AuthenticationHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := utils.GoogleOauthConfig.AuthCodeURL(OAUTH_STATE, oauth2.AccessTypeOffline)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *AuthenticationHandler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != OAUTH_STATE {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := utils.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := utils.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var gUser googleLoginPayload
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	user := User{
		FirstName:   gUser.GivenName,
		LastName:    gUser.FamilyName,
		Email:       gUser.Email,
		IsActivated: true,
		OAuth:       constants.ProviderGoogle,
	}

	existingUser, _ := h.AuthenticationStore.GetUserByEmail(user.Email)
	if existingUser == nil {
		h.AuthenticationStore.CreateUser(&user)
		mailer.SendOAuthWelcomeEmail(user.Email, user.OAuth)
	}

	accessToken, expirationTime := utils.GenerateAccessToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    accessToken,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": accessToken})
}

func (h *AuthenticationHandler) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := utils.GitHubOauthConfig.AuthCodeURL(OAUTH_STATE, oauth2.AccessTypeOffline)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (h *AuthenticationHandler) GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != OAUTH_STATE {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := utils.GitHubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := utils.GitHubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var ghUser githubLoginPayload
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	var firstName, lastName string
	if ghUser.Name != "" {
		names := strings.SplitN(ghUser.Name, " ", 2)
		firstName = names[0]
		if len(names) > 1 {
			lastName = names[1]
		}
	}

	user := User{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       ghUser.Email,
		IsActivated: true,
		OAuth:       constants.ProviderGitHub,
	}

	existingUser, _ := h.AuthenticationStore.GetUserByEmail(user.Email)
	if existingUser == nil {
		h.AuthenticationStore.CreateUser(&user)
		mailer.SendOAuthWelcomeEmail(user.Email, user.OAuth)
	}

	accessToken, expirationTime := utils.GenerateAccessToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    accessToken,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": accessToken})
}
