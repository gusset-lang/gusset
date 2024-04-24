package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneSequenceTree(t *testing.T) {
	tree := runeSequenceTree

	assert.NotNil(t, tree)
	assert.Len(t, tree, 22)

	// test single character token
	openParen, ok := tree['(']
	assert.True(t, ok, "expected tree to have key for %s", OPEN_PAREN)
	assert.Equal(t, openParen, runeTreeNode{
		t:        optionalToken(OPEN_PAREN),
		children: nil,
	})

	// test multi character token against single character token
	add, ok := tree['+']
	assert.True(t, ok, "expected tree to have key for %s", ADD)
	assert.Equal(t, optionalToken(ADD), add.t)
	assert.Len(t, add.children, 2)
	increment, ok := add.children['+']
	assert.True(t, ok, "expected %s subtree to have key for %s", ADD, ASSIGN_INC)
	assert.Equal(t, optionalToken(ASSIGN_INC), increment.t)
	assert.Len(t, increment.children, 0)

	// test multi character token with only first-rune overlap
	assign, ok := tree['=']
	assert.True(t, ok, "expected tree to have key for %s", ASSIGN)
	assert.Equal(t, optionalToken(ASSIGN), assign.t)
	assert.Len(t, add.children, 2)
	arrow, ok := assign.children['>']
	assert.True(t, ok, "expected %s subtree to have key for %s", ASSIGN, ARROW)
	assert.Equal(t, optionalToken(ARROW), arrow.t)
	assert.Len(t, arrow.children, 0)
}
