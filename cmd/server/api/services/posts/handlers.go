package posts

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type PostsHandler struct {
	PostsStore          PostsStore
	AuthenticationStore authentication.AuthenticationStore
}

var validate = validator.New()

func (h *PostsHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload postCreatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userEmail, ok := r.Context().Value(constants.EmailKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "email not found in context")
		return
	}

	user, err := h.AuthenticationStore.GetUserByEmail(userEmail)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	post := Post{
		AuthorID: user.ID,
		Author:   *user,
		Title:    payload.Title,
		Content:  payload.Content,
	}

	if err := h.PostsStore.CreatePost(&post); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	postResponse := postCreateResponse{
		ID: post.ID,
		Author: postCreateResponseAuthor{
			ID:        post.Author.ID,
			FirstName: post.Author.FirstName,
			LastName:  post.Author.LastName,
			Email:     post.Author.Email,
		},
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusCreated, postResponse)
}
