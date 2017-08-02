package roles

import "net/http"

var role = &Role{}

// Register register role with conditions
func Register(name string, fc func(request *http.Request, user interface{}) bool) {
	role.Register(name, fc)
}

// Allow allows permission mode for roles
func Allow(mode PermissionMode, roles ...string) *Permission {
	return role.Allow(mode, roles...)
}

// Deny deny permission mode for roles
func Deny(mode PermissionMode, roles ...string) *Permission {
	return role.Deny(mode, roles...)
}

// Get role defination
func Get(name string) (func(req *http.Request, user interface{}) bool, bool) {
	return role.Get(name)
}

// Remove role definition from global role instance
func Remove(name string) {
	role.Remove(name)
}

// Reset role definitions from global role instance
func Reset() {
	role.Reset()
}

// MatchedRoles return defined roles from user
func MatchedRoles(req *http.Request, user interface{}) []string {
	return role.MatchedRoles(req, user)
}

// HasRole check if current user has role
func HasRole(req *http.Request, user interface{}, roles ...string) bool {
	return role.HasRole(req, user)
}
