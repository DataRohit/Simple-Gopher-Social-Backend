package comments

import (
	"errors"
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/cmd/server/api/services/posts"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CommentsHandler struct {
	CommentsStore       CommentsStore
	PostsStore          posts.PostsStore
	AuthenticationStore authentication.AuthenticationStore
}

var validate = validator.New()

func (h *CommentsHandler) fetchCommentDetails(commentID uuid.UUID) (*commentCreateUpdateResponse, error) {
	comment, err := h.CommentsStore.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	commentLikes, _ := h.CommentsStore.GetLikesCount(comment.ID)
	commentDislikes, _ := h.CommentsStore.GetDislikesCount(comment.ID)

	post, err := h.PostsStore.GetPostByID(comment.PostID)
	if err != nil {
		return nil, err
	}

	postLikes, _ := h.PostsStore.GetLikesCount(post.ID)
	postDislikes, _ := h.PostsStore.GetDislikesCount(post.ID)

	author, err := h.AuthenticationStore.GetUserByID(comment.AuthorID.String())
	if err != nil {
		return nil, err
	}

	postAuthor, err := h.AuthenticationStore.GetUserByID(post.AuthorID.String())
	if err != nil {
		return nil, err
	}

	return &commentCreateUpdateResponse{
		ID: comment.ID,
		Author: commentCreateUpdateResponseAuthor{
			ID:        author.ID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
			Email:     author.Email,
		},
		Post: commentCreateUpdateResponsePost{
			ID: post.ID,
			Author: commentCreateUpdateResponsePostAuthor{
				ID:        postAuthor.ID,
				FirstName: postAuthor.FirstName,
				LastName:  postAuthor.LastName,
				Email:     postAuthor.Email,
			},
			Title:     post.Title,
			Content:   post.Content,
			Likes:     postLikes,
			Dislikes:  postDislikes,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
		Content:   comment.Content,
		Likes:     commentLikes,
		Dislikes:  commentDislikes,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}, nil
}

func (h *CommentsHandler) verifyCommentOwnership(commentID uuid.UUID, authUserId string) (*Comment, error) {
	comment, err := h.CommentsStore.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	commentAuthor, err := h.AuthenticationStore.GetUserByID(comment.AuthorID.String())
	if err != nil {
		return nil, err
	}

	post, err := h.PostsStore.GetPostByID(comment.PostID)
	if err != nil {
		return nil, err
	}

	comment.Author = *commentAuthor
	comment.Post = *post

	if err := utils.VerifyOwnership(comment.AuthorID.String(), authUserId); err != nil {
		return nil, err
	}

	return comment, nil
}

func (h *CommentsHandler) handleCommentLikeDislike(userUUID, commentID uuid.UUID, action string) error {
	switch action {
	case "like":
		if liked, _ := h.CommentsStore.HasUserLikedComment(userUUID, commentID); liked {
			return errors.New("you have already liked this comment")
		}
		if disliked, _ := h.CommentsStore.HasUserDislikedComment(userUUID, commentID); disliked {
			h.CommentsStore.RemoveDislikeComment(commentID, userUUID)
		}
		return h.CommentsStore.LikeComment(&CommentLike{
			UserID:    userUUID,
			CommentID: commentID,
		})

	case "unlike":
		if liked, _ := h.CommentsStore.HasUserLikedComment(userUUID, commentID); !liked {
			return errors.New("you haven't liked this comment")
		}
		return h.CommentsStore.RemoveLikeComment(userUUID, commentID)

	case "dislike":
		if disliked, _ := h.CommentsStore.HasUserDislikedComment(userUUID, commentID); disliked {
			return errors.New("you have already disliked this comment")
		}
		if liked, _ := h.CommentsStore.HasUserLikedComment(userUUID, commentID); liked {
			h.CommentsStore.RemoveLikeComment(userUUID, commentID)
		}
		return h.CommentsStore.DislikeComment(&CommentDislike{
			UserID:    userUUID,
			CommentID: commentID,
		})

	case "undislike":
		if disliked, _ := h.CommentsStore.HasUserDislikedComment(userUUID, commentID); !disliked {
			return errors.New("you haven't disliked this comment")
		}
		return h.CommentsStore.RemoveDislikeComment(userUUID, commentID)
	}

	return nil
}

func (h *CommentsHandler) GetCommentByIDHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.fetchCommentDetails(commentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, comment)
}

func (h *CommentsHandler) GetCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	limit := r.Context().Value(constants.LimitKey).(int)
	offset := r.Context().Value(constants.OffsetKey).(int)
	orderby := r.Context().Value(constants.OrderByKey).(string)
	desc := r.Context().Value(constants.DescKey).(string) == "true"

	comments, err := h.CommentsStore.GetCommentsForPost(postID, limit, offset, orderby, desc)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	commentResponses := make([]commentCreateUpdateResponse, 0)
	for _, comment := range comments {
		commentLikes, _ := h.CommentsStore.GetLikesCount(comment.ID)
		commentDislikes, _ := h.CommentsStore.GetDislikesCount(comment.ID)

		author, err := h.AuthenticationStore.GetUserByID(comment.AuthorID.String())
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		post, err := h.PostsStore.GetPostByID(comment.PostID)
		if err != nil {
			utils.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		postLikes, _ := h.PostsStore.GetLikesCount(post.ID)
		postDislikes, _ := h.PostsStore.GetDislikesCount(post.ID)

		commentResponses = append(commentResponses, commentCreateUpdateResponse{
			ID: comment.ID,
			Author: commentCreateUpdateResponseAuthor{
				ID:        author.ID,
				FirstName: author.FirstName,
				LastName:  author.LastName,
				Email:     author.Email,
			},
			Post: commentCreateUpdateResponsePost{
				ID: post.ID,
				Author: commentCreateUpdateResponsePostAuthor{
					ID:        author.ID,
					FirstName: author.FirstName,
					LastName:  author.LastName,
					Email:     author.Email,
				},
				Title:     post.Title,
				Content:   post.Content,
				Likes:     postLikes,
				Dislikes:  postDislikes,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			},
			Content:   comment.Content,
			Likes:     commentLikes,
			Dislikes:  commentDislikes,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	utils.WriteJSON(w, http.StatusOK, commentResponses)
}

func (h *CommentsHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload commentCreateUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	comment := &Comment{
		AuthorID: uuid.MustParse(authUserID),
		PostID:   postID,
		Content:  payload.Content,
	}

	if err := h.CommentsStore.CreateComment(comment); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	commentResponse, err := h.fetchCommentDetails(comment.ID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, commentResponse)
}

func (h *CommentsHandler) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload commentCreateUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	comment, err := h.verifyCommentOwnership(commentID, authUserID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	comment.Content = payload.Content

	if err := h.CommentsStore.UpdateComment(commentID, comment); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	commentResponse, err := h.fetchCommentDetails(commentID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, commentResponse)
}

func (h *CommentsHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	_, err = h.verifyCommentOwnership(commentID, authUserID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.CommentsStore.DeleteComment(commentID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *CommentsHandler) LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	if err := h.handleCommentLikeDislike(uuid.MustParse(authUserID), commentID, "like"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *CommentsHandler) UnlikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	if err := h.handleCommentLikeDislike(uuid.MustParse(authUserID), commentID, "unlike"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *CommentsHandler) DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	if err := h.handleCommentLikeDislike(uuid.MustParse(authUserID), commentID, "dislike"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *CommentsHandler) UndislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	authUserID := r.Context().Value(constants.UserIDKey).(string)

	if err := h.handleCommentLikeDislike(uuid.MustParse(authUserID), commentID, "undislike"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
