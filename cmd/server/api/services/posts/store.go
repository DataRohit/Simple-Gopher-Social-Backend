package posts

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostsStore interface {
	CreatePost(post *Post) error
	GetPostByID(postID uuid.UUID) (*Post, error)
	GetPosts() ([]Post, error)
	UpdatePost(post *Post) error
	DeletePost(postID uuid.UUID) error
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

func (s *postsStore) GetPosts() ([]Post, error) {
	var posts []Post
	if err := s.postgresDB.Find(&posts).Error; err != nil {
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