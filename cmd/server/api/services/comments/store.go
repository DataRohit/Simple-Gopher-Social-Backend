package comments

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentsStore interface {
	CreateComment(comment *Comment) error
	GetCommentByID(commentID uuid.UUID) (*Comment, error)
	GetCommentsForPost(postID uuid.UUID, limit, offset int, orderby string, desc bool) ([]Comment, error)
	UpdateComment(commentID uuid.UUID, comment *Comment) error
	DeleteComment(commentID uuid.UUID) error
	LikeComment(like *CommentLike) error
	RemoveLikeComment(userID, commentID uuid.UUID) error
	DislikeComment(dislike *CommentDislike) error
	RemoveDislikeComment(userID, commentID uuid.UUID) error
	HasUserLikedComment(userID, commentID uuid.UUID) (bool, error)
	HasUserDislikedComment(userID, commentID uuid.UUID) (bool, error)
	GetLikesCount(commentId uuid.UUID) (int64, error)
	GetDislikesCount(commentId uuid.UUID) (int64, error)
}

type commentsStore struct {
	postgresDB *gorm.DB
}

func NewCommentStore(postgresDB *gorm.DB) CommentsStore {
	return &commentsStore{
		postgresDB: postgresDB,
	}
}

func (cs *commentsStore) CreateComment(comment *Comment) error {
	return cs.postgresDB.Create(comment).Error
}

func (cs *commentsStore) GetCommentByID(commentID uuid.UUID) (*Comment, error) {
	var comment Comment
	err := cs.postgresDB.Preload("Author").Preload("Post").First(&comment, "id = ?", commentID).Error
	return &comment, err
}

func (cs *commentsStore) GetCommentsForPost(postID uuid.UUID, limit, offset int, orderby string, desc bool) ([]Comment, error) {
	var comments []Comment

	order := orderby
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}

	if err := cs.postgresDB.Preload("Author").Preload("Post").Limit(limit).Offset(offset).Order(order).Find(&comments, "post_id = ?", postID).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (cs *commentsStore) UpdateComment(commentID uuid.UUID, comment *Comment) error {
	return cs.postgresDB.Model(&Comment{}).Where("id = ?", commentID).Updates(comment).Error
}

func (cs *commentsStore) DeleteComment(commentID uuid.UUID) error {
	return cs.postgresDB.Delete(&Comment{}, "id = ?", commentID).Error
}

func (cs *commentsStore) LikeComment(like *CommentLike) error {
	return cs.postgresDB.Create(like).Error
}

func (cs *commentsStore) RemoveLikeComment(userID, commentID uuid.UUID) error {
	return cs.postgresDB.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&CommentLike{}).Error
}

func (cs *commentsStore) DislikeComment(dislike *CommentDislike) error {
	return cs.postgresDB.Create(dislike).Error
}

func (cs *commentsStore) RemoveDislikeComment(userID, commentID uuid.UUID) error {
	return cs.postgresDB.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&CommentDislike{}).Error
}

func (cs *commentsStore) HasUserLikedComment(userID, commentID uuid.UUID) (bool, error) {
	var like CommentLike
	err := cs.postgresDB.First(&like, "user_id = ? AND comment_id = ?", userID, commentID).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cs *commentsStore) HasUserDislikedComment(userID, commentID uuid.UUID) (bool, error) {
	var dislike CommentDislike
	err := cs.postgresDB.First(&dislike, "user_id = ? AND comment_id = ?", userID, commentID).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cs *commentsStore) GetLikesCount(commentId uuid.UUID) (int64, error) {
	var count int64
	err := cs.postgresDB.Model(&CommentLike{}).Where("comment_id = ?", commentId).Count(&count).Error
	return count, err
}

func (cs *commentsStore) GetDislikesCount(commentId uuid.UUID) (int64, error) {
	var count int64
	err := cs.postgresDB.Model(&CommentDislike{}).Where("comment_id = ?", commentId).Count(&count).Error
	return count, err
}
