package url_validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLValidation(t *testing.T) {
	tests := []struct {
		name      string
		link      Link
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидные URL разных типов",
			link: Link{
				Website: "https://example.com",
				Blog:    "https://blog.example.com",
				API:     "https://api.example.com/v1", // содержит 'api'
			},
			wantError: false,
		},
		{
			name: "URL с разными схемами - все валидны",
			link: Link{
				Website: "postgres://user:pass@localhost:5432/db",
				Blog:    "clickhouse://user:pass@localhost:9000",
				API:     "https://api.example.com", // содержит 'api'
			},
			wantError: false,
		},
		{
			name: "URL без схемы - ошибка",
			link: Link{
				Website: "example.com", // нет схемы
				Blog:    "https://blog.example.com",
				API:     "https://api.example.com",
			},
			wantError: true,
			errorMsg:  "Website",
		},
		{
			name: "пустой Website - ошибка (required)",
			link: Link{
				Website: "",
				Blog:    "https://blog.example.com",
				API:     "https://api.example.com",
			},
			wantError: true,
			errorMsg:  "Website",
		},
		{
			name: "Blog опциональный - пустой допустим",
			link: Link{
				Website: "https://example.com",
				Blog:    "",
				API:     "https://api.example.com",
			},
			wantError: false,
		},
		{
			name: "Blog опциональный - невалидный URL",
			link: Link{
				Website: "https://example.com",
				Blog:    "not a url", // невалидный URL
				API:     "https://api.example.com",
			},
			wantError: true,
			errorMsg:  "Blog",
		},
		{
			name: "API не содержит 'api' - ошибка",
			link: Link{
				Website: "https://example.com",
				Blog:    "https://blog.example.com",
				API:     "https://example.com/v1", // не содержит 'api'
			},
			wantError: true,
			errorMsg:  "API",
		},
		{
			name: "API валидный с contains api",
			link: Link{
				Website: "https://example.com",
				Blog:    "https://blog.example.com",
				API:     "https://api.example.com/v1",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.link.Validate()
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

func TestEdgeCases(t *testing.T) {
	t.Run("localhost URL", func(t *testing.T) {
		link := Link{
			Website: "http://localhost:8080",
			Blog:    "",
			API:     "https://api.example.com",
		}
		err := link.Validate()
		assert.NoError(t, err)
	})

	t.Run("IP as host", func(t *testing.T) {
		link := Link{
			Website: "https://192.168.1.1",
			Blog:    "",
			API:     "https://api.example.com",
		}
		err := link.Validate()
		assert.NoError(t, err)
	})

	t.Run("URL with port", func(t *testing.T) {
		link := Link{
			Website: "https://example.com:8080/path",
			Blog:    "",
			API:     "https://api.example.com",
		}
		err := link.Validate()
		assert.NoError(t, err)
	})

	t.Run("URL with query params", func(t *testing.T) {
		link := Link{
			Website: "https://example.com?param=value",
			Blog:    "",
			API:     "https://api.example.com",
		}
		err := link.Validate()
		assert.NoError(t, err)
	})

	t.Run("DSN форматы для разных БД", func(t *testing.T) {
		link := Link{
			Website: "postgres://user:pass@localhost:5432/db",
			Blog:    "mysql://user:pass@localhost:3306/db",
			API:     "https://api.example.com", // содержит 'api'
		}
		err := link.Validate()
		assert.NoError(t, err)
	})
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидные URL разных типов",
			config: Config{
				AnyURL:        "postgres://user:pass@localhost:5432/db",
				WebURL:        "https://example.com",
				DatabaseURL:   "postgres://user:pass@localhost:5432/db",
				ClickHouseURL: "clickhouse://user:pass@localhost:9000",
			},
			wantError: false,
		},
		{
			name: "AnyURL - любой валидный URL",
			config: Config{
				AnyURL:        "ftp://files.example.com",
				WebURL:        "https://example.com",
				DatabaseURL:   "postgres://user:pass@localhost:5432/db",
				ClickHouseURL: "clickhouse://user:pass@localhost:9000",
			},
			wantError: false,
		},
		{
			name: "WebURL - только http/https",
			config: Config{
				AnyURL:        "https://example.com",
				WebURL:        "ftp://example.com",
				DatabaseURL:   "postgres://user:pass@localhost:5432/db",
				ClickHouseURL: "clickhouse://user:pass@localhost:9000",
			},
			wantError: true,
			errorMsg:  "WebURL",
		},
		{
			name: "DatabaseURL - валидный DSN",
			config: Config{
				AnyURL:        "https://example.com",
				WebURL:        "https://example.com",
				DatabaseURL:   "invalid",
				ClickHouseURL: "clickhouse://user:pass@localhost:9000",
			},
			wantError: true,
			errorMsg:  "DatabaseURL",
		},
		{
			name: "ClickHouseURL - опциональный, пустой допустим",
			config: Config{
				AnyURL:        "https://example.com",
				WebURL:        "https://example.com",
				DatabaseURL:   "postgres://user:pass@localhost:5432/db",
				ClickHouseURL: "",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
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
