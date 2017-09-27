package roles

// Permissioner permissioner interface
type Permissioner interface {
	HasPermission(mode PermissionMode, roles ...interface{}) bool
}

// ConcatPermissioner concat permissioner
func ConcatPermissioner(ps ...Permissioner) Permissioner {
	return permissioners(ps)
}

type permissioners []Permissioner

// HasPermission check has permission for permissioners or not
func (ps permissioners) HasPermission(mode PermissionMode, roles ...interface{}) bool {
	for _, p := range ps {
		if !p.HasPermission(mode, roles) {
			return false
		}
	}

	return true
}
