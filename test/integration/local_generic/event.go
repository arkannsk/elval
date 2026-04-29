package local_generic

import (
	"time"

	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

//go:generate elval-gen gen -input .

type Event struct {
	// @evl:validate required
	Name string

	// @evl:validate required
	// @evl:validate after:2020-01-01
	// @evl:validate before:2030-01-01
	StartDate model.Option[time.Time]

	// @evl:validate required
	// @evl:validate after:2020-01-01
	EndDate model.Option[time.Time]
}
