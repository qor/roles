package roles

import "errors"

type PermissionMode uint32

const (
	Read PermissionMode = 1 << (32 - 1 - iota)
	Update
	Create
	Delete
	CRUD

	// for anyone
	Anyone = "*"
)

var ErrPermissionDenied = errors.New("permission denied")

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

// Concat concat two permissions into a new Permission
func (current *Permission) Concat(permission *Permission) *Permission {
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

	appendRoles(permission)
	appendRoles(current)
	return &result
}

// HasPermission check roles has permission for mode or not
func (permission *Permission) HasPermission(mode PermissionMode, roles ...string) bool {
	if len(permission.denyRoles) != 0 {
		if denyRoles := permission.denyRoles[mode]; denyRoles != nil {
			if includeRoles(denyRoles, roles) {
				return false
			}
		}
	}

	if len(permission.allowRoles) == 0 {
		return true
	} else {
		if allowRoles := permission.allowRoles[mode]; allowRoles != nil {
			if includeRoles(allowRoles, roles) {
				return true
			}
		}
	}

	return false
}

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
