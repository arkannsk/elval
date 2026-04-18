package local_generic

import (
	"fmt"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
)

type Event struct {
	// @evl:validate required
	Name string

	// @evl:validate optional, min_time:now, max_time:now+365d
	// min_time:now — дата не раньше текущего времени
	// max_time:now+365d — дата не позже чем через год
	StartDate model.Option[time.Time]

	// @evl:validate optional, after_field:StartDate
	// если указано — должно быть после StartDate (если он тоже указан)
	EndDate model.Option[time.Time]
}

func (e Event) Validate() error {
	// Name: required
	if e.Name == "" {
		return &errs.ValidationError{
			Field:   "Name",
			Rule:    errs.ErrRequired.Rule,
			Message: errs.ErrRequired.Message,
		}
	}

	// StartDate: optional + min/max time
	if startVal, ok := e.StartDate.Value(); ok {
		now := time.Now()
		yearFromNow := now.AddDate(1, 0, 0)

		if startVal.Before(now) {
			return &errs.ValidationError{
				Field:   "StartDate",
				Rule:    "min_time:now",
				Message: "start date cannot be in the past",
			}
		}
		if startVal.After(yearFromNow) {
			return &errs.ValidationError{
				Field:   "StartDate",
				Rule:    "max_time:now+365d",
				Message: "start date cannot be more than 1 year in future",
			}
		}
	}

	// EndDate: optional + after_field logic
	if endVal, ok := e.EndDate.Value(); ok {
		if startVal, ok := e.StartDate.Value(); ok {
			if !endVal.After(startVal) {
				return &errs.ValidationError{
					Field:   "EndDate",
					Rule:    "after_field:StartDate",
					Message: fmt.Sprintf("end date must be after start date, got: %v <= %v", endVal, startVal),
				}
			}
		}
	}

	return nil
}
