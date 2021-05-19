package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_FromStream(t *testing.T) {
	reader, err := os.Open(`./fixtures/field.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	field := Field{}
	field.FromStream(scanner)

	assert.EqualValues(t, 0, field.Cells[0].Index)
	assert.EqualValues(t, 3, field.Cells[0].Rich)
}

func TestField_Export(t *testing.T) {
	reader, err := os.Open(`./fixtures/field.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	field := Field{}
	field.FromStream(scanner)

	assert.Equal(t, 783, len(field.Export()))
}

func TestState_FromStream(t *testing.T) {
	reader, err := os.Open(`./fixtures/game.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	state := State{}
	state.FromStream(scanner)

	assert.EqualValues(t, 0, state.Day)
	assert.EqualValues(t, 20, state.Nutrients)

	assert.True(t, state.Players[1].IsMine)
	assert.False(t, state.Players[1].IsWaiting)
	assert.EqualValues(t, 18, state.Players[1].Sun)
	assert.EqualValues(t, 1, state.Players[1].Score)

	assert.False(t, state.Players[0].IsMine)
	assert.True(t, state.Players[0].IsWaiting)
	assert.EqualValues(t, 19, state.Players[0].Sun)
	assert.EqualValues(t, 2, state.Players[0].Score)
}
