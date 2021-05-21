package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type (
	direction int8
	index     int8
	size      int8
	num       int16
)

const (
	cellsCount = 37

	emptyCell index = -1

	dirRgt    direction = 0
	dirTopRgt direction = 1
	dirTopLft direction = 2
	dirLft    direction = 3
	dirBtmLft direction = 4
	dirBtmRgt direction = 5
)

var (
	directions = [6]direction{dirRgt, dirTopRgt, dirTopLft, dirLft, dirBtmLft, dirBtmRgt}
)

type (
	indexSlice []index
	indexMap   map[index]index
)

func (x indexSlice) Len() int           { return len(x) }
func (x indexSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x indexSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x indexSlice) Sort()              { sort.Sort(x) }

type Cell struct {
	index   index
	rich    size
	neighs1 [6]index
	neighs2 indexSlice
	neighs3 indexSlice
	vectors [6]indexSlice
}

type Field struct {
	cells [37]*Cell
}

func (f *Field) fromStream(scanner *bufio.Scanner) {
	f.cells = [37]*Cell{}

	neighsMap := make([][]indexMap, 37)

	scanner.Scan()
	for i := 0; i < cellsCount; i++ {
		scanner.Scan()

		cell := &Cell{neighs2: make(indexSlice, 0, 12), neighs3: make(indexSlice, 0, 18)}
		fmt.Sscan(
			scanner.Text(),
			&cell.index, &cell.rich,
			&cell.neighs1[0], &cell.neighs1[1], &cell.neighs1[2], &cell.neighs1[3], &cell.neighs1[4], &cell.neighs1[5],
		)

		neighsMap[cell.index] = make([]indexMap, 4)
		neighsMap[cell.index][1] = make(indexMap, 6)
		neighsMap[cell.index][2] = make(indexMap, 12)
		neighsMap[cell.index][3] = make(indexMap, 18)

		for _, neigh1 := range cell.neighs1 {
			if neigh1 == emptyCell {
				continue
			}
			neighsMap[cell.index][1][neigh1] = neigh1
		}
		f.cells[i] = cell
	}

	// count 2nd neighs
	for _, cell := range f.cells {
		for _, neigh1 := range cell.neighs1 {
			if neigh1 == emptyCell {
				continue
			}
			for _, neigh2 := range f.cells[neigh1].neighs1 {
				if neigh2 == emptyCell || neigh2 == cell.index {
					continue
				}
				// in 1
				if _, found := neighsMap[cell.index][1][neigh2]; found {
					continue
				}
				// in 2
				if _, found := neighsMap[cell.index][2][neigh2]; found {
					continue
				}
				cell.neighs2 = append(cell.neighs2, neigh2)
				neighsMap[cell.index][2][neigh2] = neigh2
			}
		}
		cell.neighs2.Sort()
	}

	// count 3nd neighs
	for _, cell := range f.cells {
		for _, neigh2 := range cell.neighs2 {
			for _, neigh3 := range f.cells[neigh2].neighs1 {
				if neigh3 == emptyCell {
					continue
				}
				// in 1
				if _, found := neighsMap[cell.index][1][neigh3]; found {
					continue
				}
				// in 2
				if _, found := neighsMap[cell.index][2][neigh3]; found {
					continue
				}
				// in 3
				if _, found := neighsMap[cell.index][3][neigh3]; found {
					continue
				}
				cell.neighs3 = append(cell.neighs3, neigh3)
				neighsMap[cell.index][3][neigh3] = neigh3
			}
		}
		cell.neighs3.Sort()
	}

	// sun vectors
	for _, cell := range f.cells {
		for _, dir := range directions {
			cell.vectors[dir] = make(indexSlice, 0, 3)
			target := cell
			length := 1
			for {
				neigh := target.neighs1[dir]
				if neigh == emptyCell || length > 3 {
					break
				}
				cell.vectors[dir] = append(cell.vectors[dir], neigh)
				target = f.cells[neigh]
				length++
			}
		}
	}
}

func (f *Field) export() string {
	var result []string

	result = append(result, `37`)
	for _, cell := range f.cells {
		neigh := fmt.Sprintf("%d %d %d %d %d %d", cell.neighs1[0], cell.neighs1[1], cell.neighs1[2], cell.neighs1[3], cell.neighs1[4], cell.neighs1[5])
		result = append(result, fmt.Sprintf("%d %d %s", cell.index, cell.rich, neigh))
	}

	return strings.Join(result, "\n")
}

type Player struct {
	sun, score        num
	isMine, isWaiting bool
}

type Tree struct {
	index          index
	size           size
	isMine, isUsed bool
}

type State struct {
	day, nutrients size
	players        [2]*Player
	trees          []*Tree
}

func (s *State) fromStream(scanner *bufio.Scanner) {
	s.players = [2]*Player{}
	s.trees = make([]*Tree, 0, 16)

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &s.day)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &s.nutrients)
	// me
	scanner.Scan()
	me := &Player{isMine: true}
	fmt.Sscan(scanner.Text(), &me.sun, &me.score)
	s.players[1] = me
	// me
	scanner.Scan()
	opp := &Player{isMine: false}
	fmt.Sscan(scanner.Text(), &opp.sun, &opp.score, &opp.isWaiting)
	s.players[0] = opp
	var num int
	// trees
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := 0; i < num; i++ {
		scanner.Scan()
		tree := &Tree{}
		fmt.Sscan(scanner.Text(), &tree.index, &tree.size, &tree.isMine, &tree.isUsed)
		s.trees = append(s.trees, tree)
	}
	// actions
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := 0; i < num; i++ {
		scanner.Scan()
	}
}

func (s *State) export() string {
	var result []string

	result = append(result, fmt.Sprintf("%d", s.day))
	result = append(result, fmt.Sprintf("%d", s.nutrients))
	result = append(result, fmt.Sprintf("%d %d", s.players[1].sun, s.players[1].score))
	result = append(result, fmt.Sprintf("%d %d %d", s.players[0].sun, s.players[0].score, boolToInt(s.players[0].isWaiting)))

	result = append(result, fmt.Sprintf("%d", len(s.trees)))
	for _, tree := range s.trees {
		result = append(result, fmt.Sprintf("%d %d %d %d", tree.index, tree.size, boolToInt(tree.isMine), boolToInt(tree.isUsed)))
	}

	return strings.Join(result, "\n")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	field := &Field{}
	field.fromStream(scanner)

	println(field.export())

	state := &State{}
	state.fromStream(scanner)
}

func l(data ...interface{}) {
	fmt.Fprintln(os.Stderr, data...)
}

func boolToInt(value bool) index {
	var result index
	if value {
		result = 1
	}
	return result
}
