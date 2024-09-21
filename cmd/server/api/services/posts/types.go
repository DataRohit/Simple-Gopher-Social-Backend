package posts

import "github.com/google/uuid"

type postCreatePayload struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type postCreateResponseAuthor struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type postCreateResponse struct {
	ID        uuid.UUID                `json:"id"`
	Author    postCreateResponseAuthor `json:"author"`
	Title     string                   `json:"title"`
	Content   string                   `json:"content"`
	CreatedAt int64                    `json:"created_at"`
	UpdatedAt int64                    `json:"updated_at"`
}
