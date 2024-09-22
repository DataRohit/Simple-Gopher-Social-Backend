package posts

import (
	"errors"
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

func (h *PostsHandler) fetchPostDetails(postID uuid.UUID) (*postCreateUpdateResponse, error) {
	post, err := h.PostsStore.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	postAuthor, err := h.AuthenticationStore.GetUserByID(post.AuthorID.String())
	if err != nil {
		return nil, err
	}

	post.Author = *postAuthor
	likes, _ := h.PostsStore.GetLikesCount(postID)
	dislikes, _ := h.PostsStore.GetDislikesCount(postID)

	return &postCreateUpdateResponse{
		ID: post.ID,
		Author: postCreateUpdateResponseAuthor{
			ID:        post.Author.ID,
			FirstName: post.Author.FirstName,
			LastName:  post.Author.LastName,
			Email:     post.Author.Email,
		},
		Title:     post.Title,
		Content:   post.Content,
		Likes:     likes,
		Dislikes:  dislikes,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (h *PostsHandler) verifyOwnership(postID uuid.UUID, authUserID string) (*Post, error) {
	post, err := h.PostsStore.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	postAuthor, err := h.AuthenticationStore.GetUserByID(post.AuthorID.String())
	if err != nil {
		return nil, err
	}

	post.Author = *postAuthor

	if err := utils.VerifyOwnership(postAuthor.ID.String(), authUserID); err != nil {
		return nil, err
	}

	return post, nil
}

func (h *PostsHandler) handleLikeDislike(userUUID, postID uuid.UUID, action string) error {
	switch action {
	case "like":
		if liked, _ := h.PostsStore.HasLiked(userUUID, postID); liked {
			return errors.New("you have already liked this post")
		}
		if disliked, _ := h.PostsStore.HasDisliked(userUUID, postID); disliked {
			h.PostsStore.RemoveDislike(userUUID, postID)
		}
		return h.PostsStore.LikePost(userUUID, postID)

	case "unlike":
		if liked, _ := h.PostsStore.HasLiked(userUUID, postID); !liked {
			return errors.New("you haven't liked this post")
		}
		return h.PostsStore.RemoveLike(userUUID, postID)

	case "dislike":
		if disliked, _ := h.PostsStore.HasDisliked(userUUID, postID); disliked {
			return errors.New("you have already disliked this post")
		}
		if liked, _ := h.PostsStore.HasLiked(userUUID, postID); liked {
			h.PostsStore.RemoveLike(userUUID, postID)
		}
		return h.PostsStore.DislikePost(userUUID, postID)

	case "undislike":
		if disliked, _ := h.PostsStore.HasDisliked(userUUID, postID); !disliked {
			return errors.New("you haven't disliked this post")
		}
		return h.PostsStore.RemoveDislike(userUUID, postID)
	}
	return nil
}

func (h *PostsHandler) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	postResponse, err := h.fetchPostDetails(postID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, postResponse)
}

func (h *PostsHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(constants.LimitKey).(int)
	offset := r.Context().Value(constants.OffsetKey).(int)
	orderby := r.Context().Value(constants.OrderByKey).(string)
	desc := r.Context().Value(constants.DescKey).(string) == "true"

	posts, err := h.PostsStore.GetPosts(limit, offset, orderby, desc)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var postResponses []postCreateUpdateResponse
	for _, post := range posts {
		response, err := h.fetchPostDetails(post.ID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		postResponses = append(postResponses, *response)
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
		utils.WriteError(w, http.StatusUnauthorized, "user ID not found in context")
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

	postResponse, _ := h.fetchPostDetails(post.ID)
	utils.WriteJSON(w, http.StatusCreated, postResponse)
}

func (h *PostsHandler) UpdatePostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	existingPost, err := h.verifyOwnership(postID, authUserID)
	if err != nil {
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

	postResponse, _ := h.fetchPostDetails(existingPost.ID)
	utils.WriteJSON(w, http.StatusOK, postResponse)
}

func (h *PostsHandler) DeletePostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	_, err = h.verifyOwnership(postID, authUserID)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	if err := h.PostsStore.DeletePost(postID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to delete post")
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *PostsHandler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	h.handleLikeDislikeRequest(w, r, "like")
}

func (h *PostsHandler) UnlikePostHandler(w http.ResponseWriter, r *http.Request) {
	h.handleLikeDislikeRequest(w, r, "unlike")
}

func (h *PostsHandler) DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	h.handleLikeDislikeRequest(w, r, "dislike")
}

func (h *PostsHandler) UndislikePostHandler(w http.ResponseWriter, r *http.Request) {
	h.handleLikeDislikeRequest(w, r, "undislike")
}

func (h *PostsHandler) handleLikeDislikeRequest(w http.ResponseWriter, r *http.Request, action string) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value(constants.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid user UUID")
		return
	}

	if err := h.handleLikeDislike(userUUID, postID, action); err != nil {
		status := http.StatusInternalServerError
		if err == errors.New("you have already liked this post") || err == errors.New("you haven't liked this post") ||
			err == errors.New("you have already disliked this post") || err == errors.New("you haven't disliked this post") {
			status = http.StatusConflict
		}
		utils.WriteError(w, status, err.Error())
		return
	}

	postResponse, _ := h.fetchPostDetails(postID)
	utils.WriteJSON(w, http.StatusOK, postResponse)
}
