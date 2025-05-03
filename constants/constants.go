package constants

type ContextKey string

const (
	CaimsContextKey    = ContextKey("claims")
	ApiKeyClientCtxKey = ContextKey("apiKeyClient")
	ApiKeyUserCtxKey   = ContextKey("apiKeyUser")
)

func (k ContextKey) String() string {
	return string(k)
}

const (
	RoleCreator = "creator"
	RolePartner = "partner"
	RoleSearch  = "search"
)

var RouteRoles = map[string][]string{
	"creators": {RoleCreator},
	"partners": {RolePartner},
	"search":   {RoleCreator, RolePartner},
}

const (
	User           = "user-"
	Formatter      = "20060102150405"
	Bearer         = "bearer"
	Authentication = "Authentication"
)
