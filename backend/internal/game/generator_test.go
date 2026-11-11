package game

import (
	"regexp"
	"testing"
)

// Calls game.Generate, testing whether the generated
// sudoku grid contains only valid values
func TestGenerateValidContent(t *testing.T) {
	want := regexp.MustCompile("^[0-9]{81}$")
	solution, given := Generate(20)
	grid, _ := newGridFromBoard(given)
	solved, _ := solve(grid)

	if !want.MatchString(solution) && solution == gridToString(&solved) {
		t.Errorf(`Generate() = %q, want match for %#q`, solution, want)
	}
}
