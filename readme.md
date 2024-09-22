# Gopher Social Backend Server

This is the backend API for **GopherSocial** built using the Chi router in Go. It provides full JWT cookie-based user authentication, supports post/comment creation with pagination, and integrates OAuth for simplified login via Google and GitHub. PostgreSQL is used for data management, and Zap for structured logging. The backend leverages GORM as the ORM adapter.

---

## Key Features

- **User Authentication**: JWT cookie-based authentication with user activation and password reset.
- **OAuth Integration**: Supports login via Google and GitHub.
- **Post & Comment Management**: CRUD operations for posts and comments with pagination, likes, and dislikes.
- **Structured Logging**: Utilizes Zap for efficient, structured logs.
- **Security & Rate Limiting**: Protects routes with rate-limiting, and supports CORS and request recovery.

---

## Packages

- **Routing**: `github.com/go-chi/chi/v5`
- **Validation**: `github.com/go-playground/validator/v10`
- **JWT Handling**: `github.com/golang-jwt/jwt/v5`
- **Logging**: `go.uber.org/zap`
- **Password Hashing**: `golang.org/x/crypto`
- **Mailing**: `gopkg.in/gomail.v2`
- **PostgreSQL Driver**: `gorm.io/driver/postgres`
- **ORM**: `gorm.io/gorm`

---

## Middlewares

1. **User Authentication**: Validates user session via JWT cookies.
2. **CORS**: Enables Cross-Origin Resource Sharing for secure API access.
3. **Logging**: Structured logging using Zap.
4. **Ordering**: Middleware to handle resource ordering for lists.
5. **Pagination**: Provides limit & offset for resource pagination.
6. **Rate Limiter**: Limits request frequency to protect against abuse.
7. **RealIP**: Extracts real IP from request headers.
8. **Recover**: Gracefully handles panics and returns 500 error.
9. **RequestID**: Attaches a unique request ID to each request for tracking.
10. **Timeout**: Configures request timeouts to prevent long-running requests.

---

## Services (Apps)

- **Health**: Health check endpoint.
- **Authentication**: Handles user registration, login, password reset, and OAuth.
- **Posts**: CRUD operations for posts and handling likes/dislikes.
- **Comments**: CRUD operations for comments, with support for likes/dislikes.

---

## API Routes

### Health Check

- `GET /health/router`

### Authentication Routes

- `POST /auth/register`: Register a new user.
- `GET /auth/activate/{token}`: Activate user account using token.
- `POST /auth/login`: Login a user.
- `POST /auth/logout`: Logout the user and invalidate JWT cookie.
- `POST /auth/forgot-password`: Request a password reset.
- `POST /auth/reset-password/{token}`: Reset the password using token.
- `GET /auth/google/login`: Redirect to Google for OAuth login.
- `GET /auth/google/callback`: Google OAuth callback.
- `GET /auth/github/login`: Redirect to GitHub for OAuth login.
- `GET /auth/github/callback`: GitHub OAuth callback.

### Post Routes

- `GET /api/v1/posts/{postID}`: Get a specific post by ID.
- `GET /api/v1/posts`: Get all posts with pagination support.
- `POST /api/v1/posts`: Create a new post.
- `PATCH /api/v1/posts/{postID}`: Update an existing post by ID.
- `DELETE /api/v1/posts/{postID}`: Delete a post by ID.
- `POST /api/v1/posts/{postID}/like`: Like a post.
- `DELETE /api/v1/posts/{postID}/like`: Remove like from a post.
- `POST /api/v1/posts/{postID}/dislike`: Dislike a post.
- `DELETE /api/v1/posts/{postID}/dislike`: Remove dislike from a post.

### Comment Routes

- `GET /api/v1/comments/{commentID}`: Get a specific comment by ID.
- `GET /api/v1/posts/{postID}/comments`: Get comments for a specific post.
- `POST /api/v1/posts/{postID}/comments`: Add a comment to a post.
- `PUT /api/v1/comments/{commentID}`: Update a comment by ID.
- `DELETE /api/v1/comments/{commentID}`: Delete a comment by ID.
- `POST /api/v1/comments/{commentID}/like`: Like a comment.
- `DELETE /api/v1/comments/{commentID}/like`: Remove like from a comment.
- `POST /api/v1/comments/{commentID}/dislike`: Dislike a comment.
- `DELETE /api/v1/comments/{commentID}/dislike`: Remove dislike from a comment.

---

## Database

- **PostgreSQL**: Uses the latest Docker image of PostgreSQL for database management. The database schema is managed through GORM migrations.
