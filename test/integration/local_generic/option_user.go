package local_generic

import (
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type UserProfile struct {
	// @evl:validate required
	// @evl:validate pattern:^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$
	Email model.Option[string]

	// @evl:validate optional
	// @evl:validate min:18
	Age model.Option[int]

	// @evl:validate required
	Metadata model.Option[UserMeta]
}

type UserMeta struct {
	// @evl:validate required
	DisplayName string
}
