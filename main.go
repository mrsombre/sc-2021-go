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

var (
	scanner = bufio.NewScanner(os.Stdin)
)

func boolToInt(value bool) int {
	var result int
	if value {
		result = 1
	}
	return result
}

type indexes []index

func (x indexes) Len() int           { return len(x) }
func (x indexes) Less(i, j int) bool { return x[i] < x[j] }
func (x indexes) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x indexes) Sort()              { sort.Sort(x) }

type Field struct {
	Cells [37]*Cell
}

type Cell struct {
	index   index
	rich    size
	neighs1 [6]index
	neighs2 indexes
	neighs3 indexes
}

func (f *Field) FromStream(scanner *bufio.Scanner) {
	f.Cells = [37]*Cell{}

	neighsMap := make(map[index]map[int]map[index]index, 37)

	scanner.Scan()
	for i := 0; i < cellsCount; i++ {
		scanner.Scan()

		cell := &Cell{neighs2: make([]index, 0, 12), neighs3: make([]index, 0, 18)}
		fmt.Sscan(
			scanner.Text(),
			&cell.index, &cell.rich,
			&cell.neighs1[0], &cell.neighs1[1], &cell.neighs1[2], &cell.neighs1[3], &cell.neighs1[4], &cell.neighs1[5],
		)

		neighsMap[cell.index] = make(map[int]map[index]index, 4)
		neighsMap[cell.index][1] = make(map[index]index, 6)
		neighsMap[cell.index][2] = make(map[index]index)
		neighsMap[cell.index][3] = make(map[index]index)

		for _, neigh1 := range cell.neighs1 {
			if neigh1 == emptyCell {
				continue
			}
			neighsMap[cell.index][1][neigh1] = neigh1
		}
		f.Cells[i] = cell
	}

	// count 2nd neighs
	for _, cell := range f.Cells {
		for _, neigh1 := range cell.neighs1 {
			if neigh1 == emptyCell {
				continue
			}
			for _, neigh2 := range f.Cells[neigh1].neighs1 {
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
	for _, cell := range f.Cells {
		for _, neigh2 := range cell.neighs2 {
			for _, neigh3 := range f.Cells[neigh2].neighs1 {
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
}

func (f *Field) Export() string {
	var result []string

	result = append(result, `37`)
	for _, cell := range f.Cells {
		neigh := fmt.Sprintf("%d %d %d %d %d %d", cell.neighs1[0], cell.neighs1[1], cell.neighs1[2], cell.neighs1[3], cell.neighs1[4], cell.neighs1[5])
		result = append(result, fmt.Sprintf("%d %d %s", cell.index, cell.rich, neigh))
	}

	return strings.Join(result, "\n")
}

type State struct {
	Day, Nutrients size
	Players        [2]*Player
	Trees          []*Tree
}

type Player struct {
	Sun, Score        num
	IsMine, IsWaiting bool
}

type Tree struct {
	Index          index
	Size           size
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
	var num int
	// trees
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := 0; i < num; i++ {
		scanner.Scan()
		tree := &Tree{}
		fmt.Sscan(scanner.Text(), &tree.Index, &tree.Size, &tree.IsMine, &tree.IsUsed)
		s.Trees = append(s.Trees, tree)
	}
	// actions
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &num)
	for i := 0; i < num; i++ {
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

func l(data ...interface{}) {
	fmt.Fprintln(os.Stderr, data...)
}
