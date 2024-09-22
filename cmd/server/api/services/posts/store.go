package posts

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostsStore interface {
	CreatePost(post *Post) error
	GetPostByID(postID uuid.UUID) (*Post, error)
	GetPosts(limit, offset int, orderby string, desc bool) ([]Post, error)
	UpdatePost(post *Post) error
	DeletePost(postID uuid.UUID) error
	LikePost(userID, postID uuid.UUID) error
	DislikePost(userID, postID uuid.UUID) error
	RemoveLike(userID, postID uuid.UUID) error
	RemoveDislike(userID, postID uuid.UUID) error
	HasLiked(userID, postID uuid.UUID) (bool, error)
	HasDisliked(userID, postID uuid.UUID) (bool, error)
	GetLikesCount(postID uuid.UUID) (int64, error)
	GetDislikesCount(postID uuid.UUID) (int64, error)
}

type postsStore struct {
	postgresDB *gorm.DB
}

func NewPostsStore(postgresDB *gorm.DB) PostsStore {
	return &postsStore{
		postgresDB: postgresDB,
	}
}

func (s *postsStore) CreatePost(post *Post) error {
	return s.postgresDB.Create(post).Error
}

func (s *postsStore) GetPostByID(postID uuid.UUID) (*Post, error) {
	var post Post
	if err := s.postgresDB.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *postsStore) GetPosts(limit, offset int, orderby string, desc bool) ([]Post, error) {
	var posts []Post

	order := orderby
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}

	if err := s.postgresDB.Limit(limit).Offset(offset).Order(order).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *postsStore) UpdatePost(post *Post) error {
	return s.postgresDB.Save(post).Error
}

func (s *postsStore) DeletePost(postID uuid.UUID) error {
	return s.postgresDB.Delete(&Post{}, postID).Error
}

func (s *postsStore) LikePost(userID, postID uuid.UUID) error {
	return s.postgresDB.Create(&PostLike{UserID: userID, PostID: postID}).Error
}

func (s *postsStore) DislikePost(userID, postID uuid.UUID) error {
	return s.postgresDB.Create(&PostDislike{UserID: userID, PostID: postID}).Error
}

func (s *postsStore) RemoveLike(userID, postID uuid.UUID) error {
	return s.postgresDB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&PostLike{}).Error
}

func (s *postsStore) RemoveDislike(userID, postID uuid.UUID) error {
	return s.postgresDB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&PostDislike{}).Error
}

func (s *postsStore) HasLiked(userID, postID uuid.UUID) (bool, error) {
	var like PostLike
	err := s.postgresDB.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *postsStore) HasDisliked(userID, postID uuid.UUID) (bool, error) {
	var dislike PostDislike
	err := s.postgresDB.Where("user_id = ? AND post_id = ?", userID, postID).First(&dislike).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *postsStore) GetLikesCount(postID uuid.UUID) (int64, error) {
	var count int64
	err := s.postgresDB.Model(&PostLike{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (s *postsStore) GetDislikesCount(postID uuid.UUID) (int64, error) {
	var count int64
	err := s.postgresDB.Model(&PostDislike{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
