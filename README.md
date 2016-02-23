# Roles

Roles is an authorization library for Golang

[![GoDoc](https://godoc.org/github.com/qor/roles?status.svg)](https://godoc.org/github.com/qor/roles)

## Usage

### Permission Modes

Roles has defined [5 permission modes](https://github.com/qor/roles/blob/master/permission.go#L8-L12):

```
roles.Read
roles.Update
roles.Create
roles.Delete
roles.CRUD   // CRUD means Read, Update, Create, Delete
```

You could use those modes to [define permission](#define-permission) or create your own mode

### Define Permission

```go
import "github.com/qor/roles"

func main() {
  // Allow Permission
  permission := roles.Allow(roles.Read, "admin") // `admin` has `Read` permission, `admin` is a role name

  // Deny Permission
  permission := roles.Deny(roles.Create, "user") // `user` has no `Create` permission

  // Using Chain
  permission := roles.Allow(roles.CRUD, "admin").Allow(roles.Read, "visitor") // `admin` has `CRUD` permissions, `visitor` only has `Read` permission
  permission := roles.Allow(roles.CRUD, "admin").Deny(roles.Update, "user") // `admin` has `CRUD` permissions, `user` doesn't has `Update` permission

  // roles `Anyone` means for anyone
  permission := roles.Deny(roles.Update, roles.Anyone) // no one has update permission
}
```

### Check Permission

```go
import "github.com/qor/roles"

func main() {
  permission := roles.Allow(roles.CRUD, "admin").Deny(roles.Create, "manager").Allow(roles.Read, "visitor")

  // check if role `admin` has the Read permission
  permission.HasPermission(roles.Read, "admin")     // => true

  // check if role `admin` has the Create permission
  permission.HasPermission(roles.Create, "admin")     // => true

  // check if role `user` has the Read permission
  permission.HasPermission(roles.Read, "user")     // => true

  // check if role `user` has the Create permission
  permission.HasPermission(roles.Create, "user")     // => false

  // check if role `visitor` has the Read permission
  permission.HasPermission(roles.Read, "user")     // => true

  // check if role `visitor` has the Create permission
  permission.HasPermission(roles.Create, "user")     // => false

  // Check with multiple roles
  // check if role `admin` or `user` has the Create permission
  permission.HasPermission(roles.Create, "admin", "user")     // => true
}
```

### Register Roles

When check permission, you need to know current user's roles first, it might be boring, especially you have defined many roles based on lots of conditions,
So Roles provides some help methods to make it easier

```go
import "github.com/qor/roles"

func main() {
  // Register roles based on some conditions
  roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
      return req.RemoteAddr == "127.0.0.1" || (currentUser != nil && currentUser.(*User).Role == "admin")
  })

  roles.Register("user", func(req *http.Request, currentUser interface{}) bool {
    return currentUser != nil
  })

  roles.Register("visitor", func(req *http.Request, currentUser interface{}) bool {
    return currentUser == nil
  })

  // Get roles from a user
  matchedRoles := roles.MatchedRoles(httpRequest, user) // []string{"user", "admin"}

  // Check if role `user` or `admin` has Read permission
  permission.HasPermission(roles.Read, matchedRoles...)
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
