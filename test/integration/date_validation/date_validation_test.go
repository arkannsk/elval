package date_validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateValidation(t *testing.T) {
	tests := []struct {
		name      string
		event     Event
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидные даты",
			event: Event{
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: "2024-01-15T10:30:00.123456789Z",
				DateOnly:  "2024-01-15",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: false,
		},
		{
			name: "RFC3339Nano формат",
			event: Event{
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: "2024-01-15T10:30:00.123456789Z",
				DateOnly:  "2024-01-15",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: false,
		},
		{
			name: "невалидный CreatedAt - не RFC3339",
			event: Event{
				CreatedAt: "2024-01-15",
				UpdatedAt: "2024-01-15T10:30:00Z",
				DateOnly:  "2024-01-15",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: true,
			errorMsg:  "CreatedAt",
		},
		{
			name: "невалидный UpdatedAt - ни один формат не подходит",
			event: Event{
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: "invalid date",
				DateOnly:  "2024-01-15",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: true,
			errorMsg:  "UpdatedAt",
		},
		{
			name: "DateOnly опциональный - пустой допустим",
			event: Event{
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: "2024-01-15T10:30:00Z",
				DateOnly:  "",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: false,
		},
		{
			name: "DateOnly опциональный - невалидный формат",
			event: Event{
				CreatedAt: "2024-01-15T10:30:00Z",
				UpdatedAt: "2024-01-15T10:30:00Z",
				DateOnly:  "2024/01/15",
				Timestamp: "2024-01-15T10:30:00",
			},
			wantError: true,
			errorMsg:  "DateOnly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Ошибка: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogEntryValidation(t *testing.T) {
	tests := []struct {
		name      string
		log       LogEntry
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидные даты с timezone",
			log: LogEntry{
				EventTime:   "2024-01-15T10:30:00+03:00",
				KitchenTime: "3:04PM",
			},
			wantError: false,
		},
		{
			name: "валидные даты RFC3339",
			log: LogEntry{
				EventTime:   "2024-01-15T10:30:00Z",
				KitchenTime: "",
			},
			wantError: false,
		},
		{
			name: "невалидный EventTime",
			log: LogEntry{
				EventTime:   "2024-01-15",
				KitchenTime: "",
			},
			wantError: true,
			errorMsg:  "EventTime",
		},
		{
			name: "KitchenTime опциональный - пустой допустим",
			log: LogEntry{
				EventTime:   "2024-01-15T10:30:00Z",
				KitchenTime: "",
			},
			wantError: false,
		},
		{
			name: "KitchenTime невалидный формат",
			log: LogEntry{
				EventTime:   "2024-01-15T10:30:00Z",
				KitchenTime: "15:04",
			},
			wantError: true,
			errorMsg:  "KitchenTime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
