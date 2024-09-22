package posts

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PostsHandler struct {
	PostsStore          PostsStore
	AuthenticationStore authentication.AuthenticationStore
}

var validate = validator.New()

func (h *PostsHandler) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.PostsStore.GetPostByID(postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	postAuthor, err := h.AuthenticationStore.GetUserByID(post.AuthorID.String())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	post.Author = *postAuthor

	postResponse := postCreateUpdateResponse{
		ID: post.ID,
		Author: postCreateUpdateResponseAuthor{
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

func (h *PostsHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostsStore.GetPosts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var postResponses []postCreateUpdateResponse
	for _, post := range posts {
		postAuthor, err := h.AuthenticationStore.GetUserByID(post.AuthorID.String())
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		post.Author = *postAuthor

		postResponse := postCreateUpdateResponse{
			ID: post.ID,
			Author: postCreateUpdateResponseAuthor{
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

		postResponses = append(postResponses, postResponse)

	}

	utils.WriteJSON(w, http.StatusOK, postResponses)
}

func (h *PostsHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload postCreateUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "email not found in context")
		return
	}

	user, err := h.AuthenticationStore.GetUserByID(userID)
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

	postResponse := postCreateUpdateResponse{
		ID: post.ID,
		Author: postCreateUpdateResponseAuthor{
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

func (h *PostsHandler) UpdatePostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingPost, err := h.PostsStore.GetPostByID(postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingPostAuthor, err := h.AuthenticationStore.GetUserByID(existingPost.AuthorID.String())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingPost.Author = *existingPostAuthor

	authUserID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "email not found in context")
		return
	}

	if err := utils.VerifyOwnership(existingPostAuthor.ID.String(), authUserID); err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var payload postCreateUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingPost.Title = payload.Title
	existingPost.Content = payload.Content

	if err := h.PostsStore.UpdatePost(existingPost); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to update post")
		return
	}

	postResponse := postCreateUpdateResponse{
		ID: existingPost.ID,
		Author: postCreateUpdateResponseAuthor{
			ID:        existingPost.Author.ID,
			FirstName: existingPost.Author.FirstName,
			LastName:  existingPost.Author.LastName,
			Email:     existingPost.Author.Email,
		},
		Title:     existingPost.Title,
		Content:   existingPost.Content,
		CreatedAt: existingPost.CreatedAt,
		UpdatedAt: existingPost.UpdatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, postResponse)
}

func (h *PostsHandler) DeletePostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingPost, err := h.PostsStore.GetPostByID(postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingPostAuthor, err := h.AuthenticationStore.GetUserByID(existingPost.AuthorID.String())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existingPost.Author = *existingPostAuthor

	authUserID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "email not found in context")
		return
	}

	if err := utils.VerifyOwnership(existingPostAuthor.ID.String(), authUserID); err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	if err := h.PostsStore.DeletePost(postID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to delete post")
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
