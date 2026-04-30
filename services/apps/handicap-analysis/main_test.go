package main

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCalculateNetScore(t *testing.T) {
	assert.Equal(t, 4, calculateNetScoreForHole(5, 10, 12))
	assert.Equal(t, 5, calculateNetScoreForHole(5, 15, 12))
	assert.Equal(t, 4, calculateNetScoreForHole(6, 1, 20))
	assert.Equal(t, 5, calculateNetScoreForHole(6, 5, 20))
	assert.Equal(t, 4, calculateNetScoreForHole(3, 18, -2))
	assert.Equal(t, 3, calculateNetScoreForHole(3, 10, -2))
}
