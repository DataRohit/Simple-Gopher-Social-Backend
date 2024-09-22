# Gopher Social Backend Server

This Go backend API for GopherSocial is built with chi router, JWT cookie user authentication, and supports post/comment creation with pagination. It uses PostgreSQL for data management, Zap for structured logging, and uses gorm as database ORM adapter. It also has full user authentication with user activation and password reset. Google OAuth & GitHub OAuth added for better & easier user authentication.

## Packages Used

1. github.com/go-chi/chi/v5
2. github.com/go-playground/validator/v10
3. github.com/golang-jwt/jwt/v5
4. go.uber.org/zap
5. golang.org/x/crypto
6. gopkg.in/gomail.v2
7. gorm.io/driver/postgres
8. gorm.io/gorm

## Middlewares Used

1. User Cookie Authentication
2. CORS
3. Logging
4. Ordering
5. Pagination (Limit & Offset)
6. Rate Limiter
7. RealIP
8. Recover
9. RequestID
10. Timeout

## Services (Apps)

1. Health
2. Authentication
3. Posts
4. Comments

## Routes

    - `/health/router` - GET

    - `/auth/register` - POST
    - `/auth/activate/{token}` - GET
    - `/auth/login` - POST
    - `/auth/logout` - POST
    - `/auth/forgot-password` - POST
    - `/auth/reset-password/{token}` - POST
    - `/auth/google/login` - GET
    - `/auth/google/callback` - GET
    - `/auth/github/login` - GET
    - `/auth/github/callback` - GET

    - `/api/v1/posts/{postID}` - GET
    - `/api/v1/posts` - GET
    - `/api/v1/posts` - POST
    - `/api/v1/posts/{postID}` - PATCH
    - `/api/v1/posts/{postID}` - DELETE
    - `/api/v1/posts/{postID}/like` - POST
    - `/api/v1/posts/{postID}/like` - DELETE
    - `/api/v1/posts/{postID}/dislike` - POST
    - `/api/v1/posts/{postID}/dislike` - DELETE

    - `/api/v1/comments/{commentID}` - GET
    - `/api/v1/posts/{postID}/comments` - GET
    - `/api/v1/posts/{postID}/comments` - POST
    - `/api/v1/comments/{commentID}` - PUT
    - `/api/v1/comments/{commentID}` - DELETE
    - `/api/v1/comments/{commentID}/like` - POST
    - `/api/v1/comments/{commentID}/like` - DELETE
    - `/api/v1/comments/{commentID}/dislike` - POST
    - `/api/v1/comments/{commentID}/dislike` - DELETE

## Database

    - PostgreSQL (Latest Docker Image)
