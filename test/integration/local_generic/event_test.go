//go:build integration

package local_generic

import (
	"context"
	"testing"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестовые данные
var (
	now        = time.Now()
	pastDate   = time.Date(2019, 6, 15, 10, 30, 0, 0, time.UTC) // до after:2020-01-01
	futureDate = time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC) // валидная дата
	lateDate   = time.Date(2035, 6, 15, 10, 30, 0, 0, time.UTC) // после before:2030-01-01
)

func TestEvent_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		event     Event
		wantErr   bool
		wantRule  string // ожидаемое правило в ошибке (если есть)
		wantField string // ожидаемое поле в ошибке (если есть)
	}{
		{
			name: "valid event with required name and valid dates",
			event: Event{
				Name:      "Conference",
				StartDate: model.Some(futureDate),
				EndDate:   model.Some(futureDate.Add(24 * time.Hour)),
			},
			wantErr: false,
		},
		{
			name: "missing required name",
			event: Event{
				Name:      "", // пусто → ошибка required
				StartDate: model.Some(futureDate),
			},
			wantErr:   true,
			wantRule:  "required",
			wantField: "Name",
		},
		{
			name: "start date in past (before after:2020-01-01)",
			event: Event{
				Name:      "Old Event",
				StartDate: model.Some(pastDate), // 2019 < 2020 → ошибка after
			},
			wantErr:   true,
			wantRule:  "after",
			wantField: "StartDate",
		},
		{
			name: "end date after before limit",
			event: Event{
				Name:    "Future Event",
				EndDate: model.Some(lateDate), // 2035 > 2030 → ошибка before
			},
			wantErr:   true,
			wantRule:  "before",
			wantField: "EndDate",
		},
		{
			name: "optional dates absent — should pass",
			event: Event{
				Name: "Minimal Event",
				// StartDate и EndDate не заданы (модель None)
			},
			wantErr: false,
		},
		{
			name: "required field missing for optional wrapper",
			event: Event{
				// Name отсутствует → ошибка required
				StartDate: model.Some(futureDate),
			},
			wantErr:   true,
			wantRule:  "required",
			wantField: "Name",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.event.Validate()

			if tt.wantErr {
				require.Error(t, err, "ожидалась ошибка валидации")

				// Проверяем структуру ошибки, если указаны ожидания
				if tt.wantRule != "" || tt.wantField != "" {
					var ve *errs.ValidationError
					require.ErrorAs(t, err, &ve, "ошибка должна быть типа *errs.ValidationError")

					if tt.wantRule != "" {
						assert.Equal(t, tt.wantRule, ve.Rule, "неверное правило валидации")
					}
					if tt.wantField != "" {
						assert.Equal(t, tt.wantField, ve.Field, "неверное поле в ошибке")
					}
				}
			} else {
				require.NoError(t, err, "не ожидалась ошибка валидации")
			}
		})
	}
}

// Тест на метод Decorate (если он есть)
func TestEvent_Decorate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	event := Event{
		Name:      "Test Event",
		StartDate: model.Some(time.Now()),
	}

	err := event.Decorate(ctx)
	require.NoError(t, err, "Decorate не должен возвращать ошибку для валидных данных")
}

// Тест на работу с model.Option: IsPresent/Value
func TestEvent_OptionBehavior(t *testing.T) {
	t.Parallel()

	t.Run("IsPresent returns false for None", func(t *testing.T) {
		var opt model.Option[time.Time] // по умолчанию — None
		assert.False(t, opt.IsPresent(), "пустой Option должен возвращать IsPresent=false")
	})

	t.Run("Value returns ok=true for Some", func(t *testing.T) {
		opt := model.Some(futureDate)
		val, ok := opt.Value()
		assert.True(t, ok, "Some должен возвращать ok=true")
		assert.Equal(t, futureDate, val, "значение должно совпадать")
	})

	t.Run("Validate skips optional field when absent", func(t *testing.T) {
		// Event с отсутствующим опциональным полем должен валидироваться успешно
		event := Event{
			Name: "Test",
			// EndDate не задан (None)
		}
		err := event.Validate()
		assert.NoError(t, err, "опциональное поле без значения не должно вызывать ошибку")
	})
}
