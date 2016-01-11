# Roles

Roles is an authorization library for Golang

When we building web applications, we found we alwalys check permission with informations from HTTP Request, like Remote IP, Header, Session or using Current User, so to save our life, we extracted the Roles package when we building [QOR](http://getqor.com)

## Usage

### Permission Modes

Roles has defined [5 permission modes](https://github.com/qor/roles/blob/master/permission.go#L7-L13):

```
roles.Read
roles.Update
roles.Create
roles.Delete
role.CRUD    // CRUD means Read, Update, Create, Delete
```

You define permissions with those modes or create your own mode

### Define Permission

```go
import "github.com/qor/roles"

func main() {
  // Allow Permission
  permission := roles.Allow(roles.Read, "admin") // `admin` is a role name

  // Deny Permission
  permission := roles.Deny(roles.Create, "user")

  // Using Chain
  permission := roles.Allow(roles.CRUD, "admin").Allow(roles.Read, "visitor")
  permission := roles.Allow(roles.CRUD, "admin").Deny(roles.Update, "user")

  // roles defined a constant `Anyone` means for anyone
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

When check permission, we need to get current roles then check if he has permission,
it might be boring, especially you have defined many roles with many conditions,
So Roles provides some help methods to let you define and get roles easier

### Register Roles

```go
import "github.com/qor/roles"

func main() {
  roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
      return req.RemoteAddr == "127.0.0.1" || (currentUser != nil && currentUser.(*User).Role == "admin")
  })

  roles.Register("user", func(req *http.Request, currentUser interface{}) bool {
    return currentUser != nil
  })

  roles.Register("visitor", func(req *http.Request, currentUser interface{}) bool {
    return currentUser == nil
  })

  // Get Matched Roles
  matchedRoles := roles.MatchedRoles(httpRequest, currentUser) // []string{"user", "admin"}

  // Check if role `user` or `admin` has Read permission
  permission.HasPermission(roles.Read, matchedRoles...)
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
