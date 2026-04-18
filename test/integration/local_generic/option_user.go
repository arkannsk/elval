package local_generic

import (
	"github.com/arkannsk/elval/pkg/errs"
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
		return &errs.ValidationError{
			Field:   "DisplayName",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}
	}
	if len(m.DisplayName) > 50 {
		return errs.NewValidationError(
			"max:50",
			"display name must not exceed %d characters, got: %d",
			50, len(m.DisplayName),
		)
	}
	return nil
}
