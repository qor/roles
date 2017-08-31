package roles

import "github.com/qor/qor"

// roler role manger interface
type roler func(data interface{}, context *qor.Context) []string
