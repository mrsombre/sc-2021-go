package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	cellsCount = 37
)

var (
	scanner = bufio.NewScanner(os.Stdin)
)

const (
	drRight       = 0
	drTopRight    = 1
	drTopLeft     = 2
	drLeft        = 3
	drBottomLeft  = 4
	drBottomRight = 5
)

var directions = [6]int{drRight, drTopRight, drTopLeft, drLeft, drBottomLeft, drBottomRight}

func l(data ...interface{}) {
	fmt.Fprintln(os.Stderr, data...)
}

func boolToInt(value bool) int {
	var result int
	if value {
		result = 1
	}
	return result
}

func inNeighs(slice []int, test int) bool {
	for _, item := range slice {
		if test == item {
			return true
		}
	}
	return false
}

type Field struct {
	Cells [37]*Cell
}

type Cell struct {
	Index   int
	Rich    uint8
	Neighs  []int
	Neighs2 []int
	Neighs3 []int
}

func (f *Field) FromStream(scanner *bufio.Scanner) {
	scanner.Scan()
	for i := 0; i < cellsCount; i++ {
		scanner.Scan()
		cell := &Cell{Neighs: make([]int, 6)}
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

	// count 2nd neighs
	for _, cell := range f.Cells {
		for _, neigh1 := range cell.Neighs {
			if neigh1 == -1 {
				continue
			}
			for _, neigh2 := range f.Cells[neigh1].Neighs {
				if neigh2 == -1 || neigh2 == cell.Index {
					continue
				}
				// check
				if inNeighs(cell.Neighs, neigh2) || inNeighs(cell.Neighs2, neigh2) {
					continue
				}
				cell.Neighs2 = append(cell.Neighs2, neigh2)
			}
		}
		sort.Ints(cell.Neighs2)
	}

	// count 3nd neighs
	for _, cell := range f.Cells {
		for _, neigh2 := range cell.Neighs2 {
			for _, neigh3 := range f.Cells[neigh2].Neighs {
				if neigh3 == -1 || neigh3 == cell.Index {
					continue
				}
				// check
				if inNeighs(cell.Neighs, neigh3) || inNeighs(cell.Neighs2, neigh3) || inNeighs(cell.Neighs3, neigh3) {
					continue
				}
				cell.Neighs3 = append(cell.Neighs3, neigh3)
			}
		}
		sort.Ints(cell.Neighs3)
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
	Index          int
	Size           uint8
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

func (s *State) Export() string {
	var result []string

	result = append(result, fmt.Sprintf("%d", s.Day))
	result = append(result, fmt.Sprintf("%d", s.Nutrients))
	result = append(result, fmt.Sprintf("%d %d", s.Players[1].Sun, s.Players[1].Score))
	result = append(result, fmt.Sprintf("%d %d %d", s.Players[0].Sun, s.Players[0].Score, boolToInt(s.Players[0].IsWaiting)))

	result = append(result, fmt.Sprintf("%d", len(s.Trees)))
	for _, tree := range s.Trees {
		result = append(result, fmt.Sprintf("%d %d %d %d", tree.Index, tree.Size, boolToInt(tree.IsMine), boolToInt(tree.IsUsed)))
	}

	return strings.Join(result, "\n")
}

func main() {
	field := &Field{}
	field.FromStream(scanner)

	println(field.Export())
}
