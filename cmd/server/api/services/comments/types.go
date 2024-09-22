package comments

import "github.com/google/uuid"

type commentCreateUpdatePayload struct {
	Content string `json:"content" validate:"required,min=1,max=10000"`
}

type commentCreateUpdateResponseAuthor struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type commentCreateUpdateResponsePostAuthor struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type commentCreateUpdateResponsePost struct {
	ID        uuid.UUID                             `json:"id"`
	Author    commentCreateUpdateResponsePostAuthor `json:"author"`
	Title     string                                `json:"title"`
	Content   string                                `json:"content"`
	Likes     int64                                 `json:"likes"`
	Dislikes  int64                                 `json:"dislikes"`
	CreatedAt int64                                 `json:"created_at"`
	UpdatedAt int64                                 `json:"updated_at"`
}

type commentCreateUpdateResponse struct {
	ID        uuid.UUID                         `json:"id"`
	Author    commentCreateUpdateResponseAuthor `json:"author"`
	Post      commentCreateUpdateResponsePost   `json:"post"`
	Content   string                            `json:"content"`
	Likes     int64                             `json:"likes"`
	Dislikes  int64                             `json:"dislikes"`
	CreatedAt int64                             `json:"created_at"`
	UpdatedAt int64                             `json:"updated_at"`
}
