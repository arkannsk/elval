package local_generic

import (
	"github.com/arkannsk/elval"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type UserProfile struct {
	// @evl:validate required, pattern:^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$
	Email model.Option[string]

	// @evl:validate optional, min:18
	Age model.Option[int]

	// @evl:validate required
	Metadata model.Option[UserMeta]
}

type UserMeta struct {
	// @evl:validate required, max:50
	DisplayName string
}

func (m UserMeta) Validate() error {
	if m.DisplayName == "" {
		return &elval.ValidationError{Field: "DisplayName", Rule: "required"}
	}
	if len(m.DisplayName) > 50 {
		return &elval.ValidationError{
			Field: "DisplayName",
			Rule:  "max:50",
		}
	}
	return nil
}
