package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	cellsCount = 37
)

var (
	scanner = bufio.NewScanner(os.Stdin)
)

const (
	drRight       uint8 = 0
	drTopRight    uint8 = 1
	drTopLeft     uint8 = 2
	drLeft        uint8 = 3
	drBottomLeft  uint8 = 4
	drBottomRight uint8 = 5
)

var directions = [...]uint8{drRight, drTopRight, drTopLeft, drLeft, drBottomLeft, drBottomRight}

func l(data ...interface{}) {
	fmt.Fprintln(os.Stderr, data...)
}

type Field struct {
	Cells [37]*Cell
}

type Cell struct {
	Index  uint8
	Rich   uint8
	Neighs [6]int8
}

func (f *Field) FromStream(scanner *bufio.Scanner) {
	scanner.Scan()
	for i := 0; i < cellsCount; i++ {
		scanner.Scan()
		cell := &Cell{}
		fmt.Sscan(
			scanner.Text(),
			&cell.Index,
			&cell.Rich,
			&cell.Neighs[drRight],
			&cell.Neighs[drTopRight],
			&cell.Neighs[drTopLeft],
			&cell.Neighs[drLeft],
			&cell.Neighs[drBottomLeft],
			&cell.Neighs[drBottomRight],
		)
		f.Cells[i] = cell
	}
}

func (f *Field) Export() string {
	var result []string
	result = append(result, `37`)
	for _, cell := range f.Cells {
		neigh := fmt.Sprintf("%d %d %d %d %d %d", cell.Neighs[0], cell.Neighs[1], cell.Neighs[2], cell.Neighs[3], cell.Neighs[4], cell.Neighs[5])
		result = append(result, fmt.Sprintf("%d %d %s", cell.Index, cell.Rich, neigh))
	}
	return strings.Join(result, "\n")
}

type State struct {
	Day, Nutrients uint8
	Players        [2]*Player
	Trees          []*Tree
}

type Player struct {
	Sun, Score        uint16
	IsMine, IsWaiting bool
}

type Tree struct {
	Index, Size    uint8
	IsMine, IsUsed bool
}

func (s *State) FromStream(scanner *bufio.Scanner) {
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &s.Day)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &s.Nutrients)
	// me
	scanner.Scan()
	me := &Player{IsMine: true}
	fmt.Sscan(scanner.Text(), &me.Sun, &me.Score)
	s.Players[1] = me
	// me
	scanner.Scan()
	opp := &Player{IsMine: false}
	fmt.Sscan(scanner.Text(), &opp.Sun, &opp.Score, &opp.IsWaiting)
	s.Players[0] = opp
	var num uint8
	// trees
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := uint8(0); i < num; i++ {
		scanner.Scan()
		tree := &Tree{}
		fmt.Sscan(scanner.Text(), &tree.Index, &tree.Size, &tree.IsMine, &tree.IsUsed)
		s.Trees = append(s.Trees, tree)
	}
	// actions
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := uint8(0); i < num; i++ {
		scanner.Scan()
	}
}

func main() {
	field := &Field{}
	field.FromStream(scanner)

	println(field.Export())
}
