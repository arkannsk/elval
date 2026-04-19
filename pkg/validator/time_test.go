package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAfter(t *testing.T) {
	rule := After("2006-01-02", "2024-01-01")

	future := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	err := rule(future)
	require.Nil(t, err)

	past := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	err = rule(past)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "after")
}

func TestBefore(t *testing.T) {
	rule := Before("2006-01-02", "2025-12-31")

	past := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	err := rule(past)
	require.Nil(t, err)

	future := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	err = rule(future)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "before")
}

func TestDurationRange(t *testing.T) {
	rule := DurationRange(1*time.Second, 10*time.Second)

	err := rule(5 * time.Second)
	require.Nil(t, err)

	err = rule(500 * time.Millisecond)
	assert.Error(t, err)

	err = rule(15 * time.Second)
	assert.Error(t, err)
}
