package posts

import "github.com/google/uuid"

type postCreateUpdatePayload struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type postCreateUpdateResponseAuthor struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type postCreateUpdateResponse struct {
	ID        uuid.UUID                      `json:"id"`
	Author    postCreateUpdateResponseAuthor `json:"author"`
	Title     string                         `json:"title"`
	Content   string                         `json:"content"`
	CreatedAt int64                          `json:"created_at"`
	UpdatedAt int64                          `json:"updated_at"`
}
