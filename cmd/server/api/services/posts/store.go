package posts

import "gorm.io/gorm"

type PostsStore interface {
	CreatePost(post *Post) error
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
