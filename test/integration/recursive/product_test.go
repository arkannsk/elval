// product_test.go
package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidNode(t *testing.T) {
	// Тестирование корректного узла
	node := &Node{
		Value: "valid",
		Next: &Node{
			Value: "next",
			Next:  nil,
		},
	}

	err := node.Validate()
	require.Nil(t, err)
}

func TestInvalidValueLength(t *testing.T) {
	// Тестирование узла с Value длиной > 50
	node := &Node{
		Value: "a very long value that definitely exceeds fifty characters (51 chars)",
		Next:  nil,
	}

	err := node.Validate()
	require.NotNil(t, err)
}

func TestInvalidNext(t *testing.T) {
	// Тестирование узла с валидным Value, но невалидным Next
	node := &Node{
		Value: "valid",
		Next: &Node{
			Value: "a very long value that definitely exceeds fifty characters (51 chars)",
			Next:  nil,
		},
	}

	if err := node.Validate(); err == nil {
		t.Fatalf("Validate() should return error for invalid Next")
	}
}

func TestDeeplyNestedNode(t *testing.T) {
	// Тестирование узла с вложенностью 5 уровней
	node := &Node{
		Value: "valid",
		Next: &Node{
			Value: "valid",
			Next: &Node{
				Value: "valid",
				Next: &Node{
					Value: "valid",
					Next: &Node{
						Value: "valid",
						Next:  nil,
					},
				},
			},
		},
	}

	err := node.Validate()
	require.Nil(t, err)
}

func TestDeeplyNestedNodeInvalidOnLvl4(t *testing.T) {
	node := &Node{
		Value: "valid",
		Next: &Node{
			Value: "valid",
			Next: &Node{
				Value: "valid",
				Next: &Node{
					Value: "a very long value that definitely exceeds fifty characters (51 chars)",
					Next: &Node{
						Value: "valid",
						Next:  nil,
					},
				},
			},
		},
	}

	err := node.Validate()
	require.NotNil(t, err)
}
