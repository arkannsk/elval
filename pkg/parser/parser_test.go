package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseComments(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     []string
	}{
		{
			name:     "parse simple comments",
			filename: "../testdata/sample.go",
			want: []string{
				"@evl:validate required",
				"@evl:validate min-max:3,50",
				"@evl:validate required",
				"@evl:validate pattern:email",
				"@evl:validate required",
				"@evl:validate min-max:1,100",
				"@evl:validate optional",
				"@evl:validate pattern:uuid",
				"@evl:validate required",
				"@evl:validate min-max:0,999.99",
				"@evl:validate min-max:1,1000",
				`@evl:validate pattern:^[A-Z]{3}-\d+$`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseComments(tt.filename)

			require.NoError(t, err)

			t.Logf("annotations list")
			for i, g := range got {
				t.Logf("  [%d] %q", i, g)
			}

			assert.Len(t, got, len(tt.want))
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestParseCommentsDetailed(t *testing.T) {
	filename := "../testdata/sample.go"
	got, err := ParseComments(filename)

	require.NoError(t, err)
	require.Len(t, got, 12)

	t.Run("Product with regexp", func(t *testing.T) {
		expected := `@evl:validate pattern:^[A-Z]{3}-\d+$` // raw string
		assert.Contains(t, got, expected)
	})
}
