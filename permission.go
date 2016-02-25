package roles

import "errors"

// PermissionMode permission mode
type PermissionMode uint32

const (
	// Create predefined permission mode, create permission
	Create PermissionMode = 1 << (32 - 1 - iota)
	// Read predefined permission mode, read permission
	Read
	// Update predefined permission mode, update permission
	Update
	// Delete predefined permission mode, deleted permission
	Delete
	// CRUD predefined permission mode, create+read+update+delete permission
	CRUD
)

// ErrPermissionDenied no permission error
var ErrPermissionDenied = errors.New("permission denied")

// NewPermission initialize a new permission for default role
func NewPermission() *Permission {
	return role.newPermission()
}

// Permission a struct contains permission definitions
type Permission struct {
	Role       *Role
	allowRoles map[PermissionMode][]string
	denyRoles  map[PermissionMode][]string
}

func includeRoles(roles []string, values []string) bool {
	for _, role := range roles {
		if role == Anyone {
			return true
		}

		for _, value := range values {
			if value == role {
				return true
			}
		}
	}
	return false
}

// Concat concat two permissions into a new one
func (permission *Permission) Concat(newPermission *Permission) *Permission {
	var result = Permission{
		Role:       role,
		allowRoles: map[PermissionMode][]string{},
		denyRoles:  map[PermissionMode][]string{},
	}

	var appendRoles = func(p *Permission) {
		if p != nil {
			result.Role = p.Role

			for mode, roles := range p.denyRoles {
				result.denyRoles[mode] = append(result.denyRoles[mode], roles...)
			}

			for mode, roles := range p.allowRoles {
				result.allowRoles[mode] = append(result.allowRoles[mode], roles...)
			}
		}
	}

	appendRoles(newPermission)
	appendRoles(permission)
	return &result
}

// HasPermission check roles has permission for mode or not
func (permission Permission) HasPermission(mode PermissionMode, roles ...string) bool {
	if len(permission.denyRoles) != 0 {
		if denyRoles := permission.denyRoles[mode]; denyRoles != nil {
			if includeRoles(denyRoles, roles) {
				return false
			}
		}
	}

	// return true if haven't define allowed roles
	if len(permission.allowRoles) == 0 {
		return true
	}

	if allowRoles := permission.allowRoles[mode]; allowRoles != nil {
		if includeRoles(allowRoles, roles) {
			return true
		}
	}

	return false
}

// Allow allows permission mode for roles
func (permission *Permission) Allow(mode PermissionMode, roles ...string) *Permission {
	if mode == CRUD {
		return permission.Allow(Create, roles...).Allow(Update, roles...).Allow(Read, roles...).Allow(Delete, roles...)
	}

	if permission.allowRoles[mode] == nil {
		permission.allowRoles[mode] = []string{}
	}
	permission.allowRoles[mode] = append(permission.allowRoles[mode], roles...)
	return permission
}

// Deny deny permission mode for roles
func (permission *Permission) Deny(mode PermissionMode, roles ...string) *Permission {
	if mode == CRUD {
		return permission.Deny(Create, roles...).Deny(Update, roles...).Deny(Read, roles...).Deny(Delete, roles...)
	}

	if permission.denyRoles[mode] == nil {
		permission.denyRoles[mode] = []string{}
	}
	permission.denyRoles[mode] = append(permission.denyRoles[mode], roles...)
	return permission
}
