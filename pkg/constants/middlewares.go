package constants

type contextKey string

const (
	UserIDKey  contextKey = "userID"
	LimitKey   contextKey = "limit"
	OffsetKey  contextKey = "offset"
	OrderByKey contextKey = "orderby"
	DescKey    contextKey = "desc"
)

const (
	DefaultLimit = 10
	MinLimit     = 1
	MaxLimit     = 20
)

const (
	DefaultOffset = 0
)

const (
	DefaultOrderBy = "created_at"
)

const (
	DefaultDesc = "true"
)

var (
	SQLKeywords = map[string]struct{}{
		"SELECT": {}, "INSERT": {}, "UPDATE": {}, "DELETE": {},
		"FROM": {}, "WHERE": {}, "JOIN": {}, "INNER": {},
		"LEFT": {}, "RIGHT": {}, "FULL": {}, "GROUP": {},
		"ORDER": {}, "BY": {}, "HAVING": {}, "LIMIT": {},
		"OFFSET": {}, "DISTINCT": {}, "CREATE": {}, "DROP": {},
		"ALTER": {}, "TABLE": {}, "INDEX": {}, "VIEW": {},
		"SET": {}, "VALUES": {}, "INTO": {}, "AS": {},
		"AND": {}, "OR": {}, "NOT": {}, "LIKE": {},
		"IS": {}, "NULL": {}, "BETWEEN": {}, "EXISTS": {},
		"UNION": {},
	}
)
