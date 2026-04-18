package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	rule := Eq(10)

	require.Nil(t, rule(10))
	assert.Error(t, rule(5))
}

func TestNeq(t *testing.T) {
	rule := Neq(0)

	require.Nil(t, rule(5))
	assert.Error(t, rule(0))
}

func TestLt(t *testing.T) {
	rule := Lt(10)

	require.Nil(t, rule(5))
	assert.Error(t, rule(10))
	assert.Error(t, rule(15))
}

func TestGt(t *testing.T) {
	rule := Gt(0)

	require.Nil(t, rule(5))
	assert.Error(t, rule(0))
	assert.Error(t, rule(-1))
}

func TestLte(t *testing.T) {
	rule := Lte(10)

	require.Nil(t, rule(5))
	require.Nil(t, rule(10))
	assert.Error(t, rule(15))
}

func TestGte(t *testing.T) {
	rule := Gte(18)

	require.Nil(t, rule(18))
	require.Nil(t, rule(25))
	assert.Error(t, rule(16))
}
