package game

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type cellOptions [9]bool
type grid [9][9]cellOptions

func newGrid() grid {
	var result grid
	for row := range result {
		for col := range result[row] {
			for value := range result[row][col] {
				result[row][col][value] = true
			}
		}
	}
	return result
}

func (options *cellOptions) toValues() []int {
	var result []int
	for value, possible := range options {
		if possible {
			result = append(result, value)
		}
	}
	return result
}

func peers(row, col int) [][2]int {
	result := make([][2]int, 0, 20)
	for i := 0; i < 9; i++ {
		if i != col {
			result = append(result, [2]int{row, i})
		}
		if i != row {
			result = append(result, [2]int{i, col})
		}
	}
	boxRow, boxCol := row/3*3, col/3*3
	for r := boxRow; r < boxRow+3; r++ {
		for c := boxCol; c < boxCol+3; c++ {
			if r != row && c != col { // the rest are already in the list
				result = append(result, [2]int{r, c})
			}
		}
	}
	return result
}

// Eliminate removes one option from a cell. If that leaves a single option,
// it removes that value from the cell's peers too, recursively.
// Returns false on a cell with no options left.
func (g *grid) eliminate(row, col, value int) bool {
	if !g[row][col][value] {
		return true // already gone
	}
	g[row][col][value] = false

	remaining := g[row][col].toValues()
	switch len(remaining) {
	case 0:
		return false
	case 1:
		for _, peer := range peers(row, col) {
			if !g.eliminate(peer[0], peer[1], remaining[0]) {
				return false
			}
		}
	}
	return true
}

// Assign collapses a cell to one value by eliminating all the others.
// Returns false when propagated elimination stumbles upon a cell with zero
// options.
func (g *grid) assign(row, col, value int) bool {
	for other := 0; other < 9; other++ {
		if other != value && !g.eliminate(row, col, other) {
			return false
		}
	}
	return true
}

// lowestEntropyCell returns the unresolved cell with the fewest options.
func (g *grid) lowestEntropyCell() (row, col int, found bool) {
	bestEntropy := 10
	for candidateRow := 0; candidateRow < 9; candidateRow++ {
		for candidateCol := 0; candidateCol < 9; candidateCol++ {
			entropy := len(g[candidateRow][candidateCol].toValues())
			if entropy > 1 && entropy < bestEntropy {
				row, col, bestEntropy, found = candidateRow, candidateCol, entropy, true
			}
		}
	}
	return row, col, found
}

// countSolutions counts completions of the grid, stopping once it has found `limit`.
func countSolutions(current grid, limit int) int {
	row, col, found := current.lowestEntropyCell()
	if !found {
		return 1 // everything is resolved: exactly one completion
	}

	total := 0
	for _, value := range current[row][col].toValues() {
		attempt := current
		if !attempt.assign(row, col, value) {
			continue
		}
		total += countSolutions(attempt, limit-total)
		if total >= limit {
			break
		}
	}
	return total
}

// Solve the grid into a valid, solved sudoku grid.
func solve(current grid) (grid, bool) {
	// 1. Find the unresolved cell with the lowest entropy (first occurrence).
	bestRow, bestCol, found := current.lowestEntropyCell()
	if !found {
		return current, true // every cell is down to one option
	}

	// 2. Shuffle the possible assignment candidates
	candidates := current[bestRow][bestCol].toValues()
	rand.Shuffle(len(candidates), func(a, b int) {
		candidates[a], candidates[b] = candidates[b], candidates[a]
	})
	for _, value := range candidates {
		attempt := current // array copy: the snapshot for backtracking
		// 3. Attempt to assign the value
		if !attempt.assign(bestRow, bestCol, value) {
			continue
		}
		// 4. recursively solve for the grid with the value assigned
		if solved, ok := solve(attempt); ok {
			return solved, true
		}
	}
	return grid{}, false
}

// newGridFromBoard builds a grid from an 81-character board ('0' = empty).
// It returns false if the givens contradict each other.
func newGridFromBoard(board string) (grid, bool) {
	puzzle := newGrid()
	for index, char := range board {
		if char == '0' {
			continue
		}
		if !puzzle.assign(index/9, index%9, int(char-'1')) {
			return grid{}, false
		}
	}
	return puzzle, true
}

func hasUniqueSolution(board string) bool {
	puzzle, ok := newGridFromBoard(board)
	return ok && countSolutions(puzzle, 2) == 1
}

// removeClues blanks cells of a solved board while the solution stays unique.
// It stops at targetClues, or earlier if no more cells can be removed.
func removeClues(solution string, targetClues int) string {
	board := []byte(solution)
	clues := 81

	for _, index := range rand.Perm(81) {
		if clues <= targetClues {
			break
		}
		removed := board[index]
		board[index] = '0'
		if hasUniqueSolution(string(board)) {
			clues--
		} else {
			board[index] = removed // this clue is needed, put it back
		}
	}
	return string(board)
}

func Generate(targetClues int) (solution string, given string) {
	grid := newGrid()
	solved, _ := solve(grid)
	result := gridToString(&solved)
	givens := removeClues(result, targetClues)
	return result, givens
}

// Builds a string representation of the grid. Assumes the grid is reduced to a
// single option per cell, as such it should only realy be used on fully
// generated grids
func gridToString(grid *grid) string {
	outBuilder := strings.Builder{}
	// Transform into a string representation
	for rowIdx := range len(grid) {
		for colIdx := range len(grid[rowIdx]) {
			for val := 0; val < len(grid[rowIdx][colIdx]); val++ {
				if grid[rowIdx][colIdx][val] {
					outBuilder.WriteString(strconv.Itoa(val + 1))
					break
				}
			}
		}
	}

	return outBuilder.String()
}

const boardSeparator = "+-------+-------+-------+\n"

// renderBoard draws a 9x9 board with lines between the 3x3 boxes.
// symbolAt returns the character to show for a cell.
func renderBoard(symbolAt func(row, col int) byte) string {
	var builder strings.Builder
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			builder.WriteString(boardSeparator)
		}
		for col := 0; col < 9; col++ {
			if col%3 == 0 {
				builder.WriteString("| ")
			}
			builder.WriteByte(symbolAt(row, col))
			builder.WriteByte(' ')
		}
		builder.WriteString("|\n")
	}
	builder.WriteString(boardSeparator)
	return builder.String()
}

// symbol is the character for one cell of the wave-function grid:
// the digit if resolved, '.' if undecided, 'X' if it has no options left.
func (g *grid) symbol(row, col int) byte {
	remaining := g[row][col].toValues()
	switch len(remaining) {
	case 0:
		return 'X'
	case 1:
		return byte('1' + remaining[0])
	default:
		return '.'
	}
}
