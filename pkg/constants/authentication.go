package constants

type UserRole string
type OAuthProvider string

const (
	RoleUser  UserRole = "user"
	RoleStaff UserRole = "staff"
	RoleAdmin UserRole = "admin"
)

const (
	ProviderNone   OAuthProvider = "none"
	ProviderGoogle OAuthProvider = "google"
	ProviderGitHub OAuthProvider = "github"
)
