//go:build integration

package local_generic

import (
	"testing"
	"time"

	"github.com/arkannsk/elval/pkg/errs"
	"github.com/arkannsk/elval/test/integration/local_generic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSettings_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		settings UserSettings
		wantErr  bool
		wantRule string
	}{
		{
			name: "valid settings with required theme",
			settings: UserSettings{
				Theme:      "dark", // валидное значение из enum
				MaxRetries: 5,
				Timeout:    model.Some(30 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "invalid theme value",
			settings: UserSettings{
				Theme:      "invalid", // не в enum [light,dark,system]
				MaxRetries: 3,
			},
			wantErr:  true,
			wantRule: "enum",
		},
		{
			name: "max retries below min",
			settings: UserSettings{
				Theme:      "light",
				MaxRetries: 0, // < min:1
			},
			wantErr:  true,
			wantRule: "min",
		},
		{
			name: "optional timeout absent — should pass",
			settings: UserSettings{
				Theme:      "system",
				MaxRetries: 10,
				// Timeout не задан
			},
			wantErr: false,
		},
		{
			name: "timeout value below min duration",
			settings: UserSettings{
				Theme:      "dark",
				Timeout:    model.Some(500 * time.Millisecond), // < min:1s
				MaxRetries: 3,
			},
			wantErr:  true,
			wantRule: "min",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.settings.Validate()

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantRule != "" {
					var ve *errs.ValidationError
					require.ErrorAs(t, err, &ve)
					assert.Equal(t, tt.wantRule, ve.Rule)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
