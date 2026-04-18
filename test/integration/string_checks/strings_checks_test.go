package string_checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainsValidation(t *testing.T) {
	tests := []struct {
		name      string
		doc       Document
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный документ",
			doc: Document{
				Name:      "README.md",
				URL:       "https://example.com",
				Content:   "hello world",
				ImageName: "img_001",
			},
			wantError: false,
		},
		{
			name: "имя не содержит README",
			doc: Document{
				Name:      "doc.md",
				URL:       "https://example.com",
				Content:   "hello world",
				ImageName: "img_001",
			},
			wantError: true,
			errorMsg:  "Name",
		},
		{
			name: "URL не начинается с https://",
			doc: Document{
				Name:      "README.md",
				URL:       "http://example.com",
				Content:   "hello world",
				ImageName: "img_001",
			},
			wantError: true,
			errorMsg:  "URL",
		},
		{
			name: "URL не заканчивается на .com",
			doc: Document{
				Name:      "README.md",
				URL:       "https://example.org",
				Content:   "hello world",
				ImageName: "img_001",
			},
			wantError: true,
			errorMsg:  "URL",
		},
		{
			name: "контент не содержит world",
			doc: Document{
				Name:      "README.md",
				URL:       "https://example.com",
				Content:   "hello go",
				ImageName: "img_001",
			},
			wantError: true,
			errorMsg:  "Content",
		},
		{
			name: "опциональное поле - пустое",
			doc: Document{
				Name:      "README.md",
				URL:       "https://example.com",
				Content:   "hello world",
				ImageName: "",
			},
			wantError: false,
		},
		{
			name: "опциональное поле - не начинается с img_",
			doc: Document{
				Name:      "README.md",
				URL:       "https://example.com",
				Content:   "hello world",
				ImageName: "photo_001",
			},
			wantError: true,
			errorMsg:  "ImageName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				t.Logf("Ошибка: %v", err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestFileValidation(t *testing.T) {
	tests := []struct {
		name      string
		file      File
		wantError bool
		errorMsg  string
	}{
		{
			name: "валидный файл",
			file: File{
				Path: "main.go",
				Name: "test_file",
			},
			wantError: false,
		},
		{
			name: "путь не заканчивается на .go",
			file: File{
				Path: "main.py",
				Name: "test_file",
			},
			wantError: true,
			errorMsg:  "Path",
		},
		{
			name: "имя не содержит test",
			file: File{
				Path: "main.go",
				Name: "file",
			},
			wantError: true,
			errorMsg:  "Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.file.Validate()
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
