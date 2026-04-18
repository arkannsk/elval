package local_generic

import (
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type Review struct {
	// @evl:validate required, min:3, max:500
	Comment string
	// @evl:validate optional, min:1, max:5
	Rating model.Option[int]
}

type Product struct {
	// @evl:validate required
	Name string

	// @evl:validate optional, not-empty
	// каждый присутствующий Review должен быть валидным
	Reviews []model.Option[Review]
}
