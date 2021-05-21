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
	field.fromStream(scanner)

	assert.Equal(t, index(7), field.cells[7].index)
	assert.EqualValues(t, index(2), field.cells[7].rich)

	assert.Equal(t, [6]index{-1, 35, 17, 33, -1, -1}, field.cells[34].neighs1)
	assert.Equal(t, indexSlice{5, 6, 14, 17, 30, 34}, field.cells[32].neighs2)
	assert.Equal(t, indexSlice{0, 3, 5, 11, 15, 25, 31}, field.cells[28].neighs3)

	assert.Equal(t, indexSlice{3, 4, 14}, field.cells[10].vectors[dirBtmLft])
	assert.Equal(t, indexSlice{18, 36}, field.cells[6].vectors[dirRgt])
	assert.Equal(t, indexSlice{25}, field.cells[11].vectors[dirTopLft])
	assert.Equal(t, indexSlice{}, field.cells[29].vectors[dirLft])
}

func TestField_Export(t *testing.T) {
	reader, err := os.Open(`./fixtures/field.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	field := Field{}
	field.fromStream(scanner)

	assert.Equal(t, 783, len(field.export()))
}

func TestState_FromStream(t *testing.T) {
	reader, err := os.Open(`./fixtures/game.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	state := State{}
	state.fromStream(scanner)

	assert.Equal(t, size(0), state.day)
	assert.Equal(t, size(20), state.nutrients)

	assert.True(t, state.players[1].isMine)
	assert.False(t, state.players[1].isWaiting)
	assert.Equal(t, num(18), state.players[1].sun)
	assert.Equal(t, num(1), state.players[1].score)

	assert.False(t, state.players[0].isMine)
	assert.True(t, state.players[0].isWaiting)
	assert.EqualValues(t, num(19), state.players[0].sun)
	assert.EqualValues(t, num(2), state.players[0].score)
}

func TestState_Export(t *testing.T) {
	reader, err := os.Open(`./fixtures/game.txt`)
	if err != nil {
		assert.NoError(t, err)
	}
	scanner := bufio.NewScanner(reader)

	state := State{}
	state.fromStream(scanner)

	assert.Equal(t, 126, len(state.export()))
}
