package show

import "github.com/gookit/gcli/v3/show/showcom"

// PrettyJSON struct
type PrettyJSON struct {
	showcom.Base
}

// NewPrettyJSON instance
func NewPrettyJSON() *PrettyJSON {
	return &PrettyJSON{}
}
