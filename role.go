package roles

import (
	"fmt"
	"net/http"
)

const (
	// Anyone is a role for any one
	Anyone = "*"
)

// New initialize a new `Role`
func New() *Role {
	return &Role{}
}

// Role is a struct contains all roles definitions
type Role struct {
	definitions map[string]func(request *http.Request, user interface{}) bool
}

// Register register role with conditions
func (role *Role) Register(name string, fc func(req *http.Request, currentUser interface{}) bool) {
	if role.definitions == nil {
		role.definitions = map[string]func(req *http.Request, currentUser interface{}) bool{}
	}

	definition := role.definitions[name]
	if definition != nil {
		fmt.Printf("%v already defined, overwrited it!\n", name)
	}
	role.definitions[name] = fc
}

// NewPermission initialize permission
func (role *Role) NewPermission() *Permission {
	return &Permission{
		Role:         role,
		AllowedRoles: map[PermissionMode][]string{},
		DeniedRoles:  map[PermissionMode][]string{},
	}
}

// Allow allows permission mode for roles
func (role *Role) Allow(mode PermissionMode, roles ...string) *Permission {
	return role.NewPermission().Allow(mode, roles...)
}

// Deny deny permission mode for roles
func (role *Role) Deny(mode PermissionMode, roles ...string) *Permission {
	return role.NewPermission().Deny(mode, roles...)
}

// Get role defination
func (role *Role) Get(name string) (func(req *http.Request, currentUser interface{}) bool, bool) {
	fc, ok := role.definitions[name]
	return fc, ok
}

// Remove role definition
func (role *Role) Remove(name string) {
	delete(role.definitions, name)
}

// Reset role definitions
func (role *Role) Reset() {
	role.definitions = map[string]func(req *http.Request, currentUser interface{}) bool{}
}

// MatchedRoles return defined roles from user
func (role *Role) MatchedRoles(req *http.Request, currentUser interface{}) (roles []string) {
	if definitions := role.definitions; definitions != nil {
		for name, definition := range definitions {
			if definition(req, currentUser) {
				roles = append(roles, name)
			}
		}
	}
	return
}
