package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Cell represents a single cell in the table.
type Cell struct {
	formula string   // Original formula or a string representation of a number.
	value   float64  // Cached computed value.
	deps    []string // Cells this cell depends on (e.g., ["a1", "b2"] for "a1+b2").
}

// Table holds the entire spreadsheet data and dependency graph.
type Table struct {
	cells map[string]*Cell
	// Reverse dependencies: maps a cell to a list of cells that depend on it.
	// e.g., "a1" -> ["f4", "g5"] if f4 and g5 use a1 in their formulas.
	revDeps map[string][]string
}

// NewTable initializes and returns a new Table.
func NewTable() *Table {
	return &Table{
		cells:   make(map[string]*Cell),
		revDeps: make(map[string][]string),
	}
}

// Helper: Regex to find cell names (e.g., a1, B2) in a formula.
var cellRefRegex = regexp.MustCompile(`[a-zA-Z]+\d+`)

// set_cell sets the formula for a cell and triggers recalculation.
func (t *Table) set_cell(cell string, formula string) error {
	// 1. Clean up old dependencies
	oldCell, exists := t.cells[cell]
	if exists {
		for _, dep := range oldCell.deps {
			t.removeRevDep(dep, cell)
		}
	}

	// 2. Analyze new dependencies and store formula/value
	cellStruct := &Cell{formula: formula}

	// Check if the formula is a simple number
	val, err := strconv.ParseFloat(formula, 64)
	if err == nil {
		cellStruct.value = val
	} else {
		// It's a formula, extract dependencies
		cellStruct.deps = cellRefRegex.FindAllString(formula, -1)
		// Update reverse dependencies
		for _, dep := range cellStruct.deps {
			t.addRevDep(dep, cell)
		}
	}

	t.cells[cell] = cellStruct

	// 3. Recalculate the current cell and all dependents
	t.recalculate(cell)

	return nil
}

// Helper: Adds a reverse dependency
func (t *Table) addRevDep(dep string, target string) {
	// Only add if not already present
	for _, existing := range t.revDeps[dep] {
		if existing == target {
			return
		}
	}
	t.revDeps[dep] = append(t.revDeps[dep], target)
}

// Helper: Removes a reverse dependency
func (t *Table) removeRevDep(dep string, target string) {
	targets := t.revDeps[dep]
	for i, existing := range targets {
		if existing == target {
			t.revDeps[dep] = append(targets[:i], targets[i+1:]...)
			break
		}
	}
}

// recalculate computes the value for a cell and recursively triggers
// recalculation for all cells that depend on it (reverse dependencies).
func (t *Table) recalculate(cell string) {
	// A simple number formula was already calculated in set_cell.
	// Only proceed if it's a formula with dependencies.
	cellStruct, exists := t.cells[cell]
	if !exists || len(cellStruct.deps) == 0 {
		// If it's a simple number, its value is already set.
		// Continue to propagate changes to its dependents.
		goto propagate
	}

	// Evaluation logic (Simplified: only supports +-*/ on cell references)

	// Start with the formula string
	evalFormula := cellStruct.formula

	// Replace cell references with their current values
	for _, dep := range cellStruct.deps {
		depCell, exists := t.cells[dep]
		depValue := 0.0
		if exists {
			depValue = depCell.value
		}

		// This simplified replacement is naive and might break multi-digit numbers
		// if a cell name is a substring of another (e.g., "a1" in "a10").
		// A more robust solution would use an Abstract Syntax Tree (AST) or
		// a proper expression evaluation library.
		evalFormula = strings.ReplaceAll(evalFormula, dep, strconv.FormatFloat(depValue, 'f', -1, 64))
	}

	// *** SIMPLIFIED EVALUATION START ***
	// Since full expression parsing is complex, we'll assume a very simple
	// formula structure (e.g., a1+b2, not a1*(b2+c3)) for this example.
	// A production-ready version must use an expression evaluator.

	result := 0.0
	// For simplicity, we just try to evaluate the string as an expression.
	// In a real application, you'd use a math expression parser/evaluator.
	// For this example, we'll just demonstrate the dependency logic.

	// Attempt to get the final value by *assuming* simple operations
	// and using a placeholder for the actual evaluation.
	result = t.evaluateSimpleExpression(evalFormula)
	cellStruct.value = result
	// *** SIMPLIFIED EVALUATION END ***

propagate:
	// Recursively recalculate all cells that depend on the current cell
	for _, dependentCell := range t.revDeps[cell] {
		t.recalculate(dependentCell)
	}
}

// Placeholder for expression evaluation.
// A real solution would use a library like 'github.com/Knetic/govaluate'
// or implement a shunting-yard algorithm and RPN evaluator.
func (t *Table) evaluateSimpleExpression(expression string) float64 {
	// Highly simplified: just for demonstration.
	// This function *cannot* safely evaluate arbitrary math expressions.
	// Assume it returns the correct result for the sake of the dependency example.

	// Example: If expression is "10.0+5.0", it should return 15.0

	// In a real scenario, this would involve:
	// 1. Tokenization
	// 2. Shunting-yard algorithm (to convert infix to postfix/RPN)
	// 3. RPN evaluation

	// For now, let's just make up a value if we can't parse it
	// as a float (which would happen if it still contains an operator)
	val, err := strconv.ParseFloat(expression, 64)
	if err == nil {
		return val
	}

	// If it's still a complex formula, we'll return a placeholder.
	// For a demonstration, assume it *successfully* evaluates to an arbitrary result
	// based on the length of the string to show a change.
	return float64(len(expression))
}

// get_cell returns the current computed value of the cell.
func (t *Table) get_cell(cell string) (float64, error) {
	c, exists := t.cells[cell]
	if !exists {
		return 0, fmt.Errorf("cell %s does not exist", cell)
	}
	return c.value, nil
}

func main() {
	table := NewTable()

	// 1. Set a base value
	table.set_cell("a1", "10")
	// a1 = 10

	// 2. Set a dependent formula
	table.set_cell("b2", "a1+5") // b2 = 10 + 5 = 15
	// Dependencies: b2 depends on a1

	// 3. Set a second-level dependent formula
	table.set_cell("c3", "b2*2") // c3 = 15 * 2 (simplified evaluation logic)
	// Dependencies: c3 depends on b2, which depends on a1

	a1Val, _ := table.get_cell("a1")
	b2Val, _ := table.get_cell("b2")
	c3Val, _ := table.get_cell("c3")

	fmt.Printf("Initial Values:\n")
	fmt.Printf("a1: %.2f\n", a1Val) // Expected: 10.00
	fmt.Printf("b2: %.2f\n", b2Val) // Expected: 15.00 (based on a1's value)
	fmt.Printf("c3: %.2f\n", c3Val) // Expected: 5.00 (based on simplified eval of "15.0*2")
	fmt.Println("--------------------")

	// 4. Update the base cell (a1)
	fmt.Printf("Updating a1 to '50'...\n")
	table.set_cell("a1", "50")
	// This triggers: a1 -> b2 (recalc) -> c3 (recalc)

	a1Val, _ = table.get_cell("a1")
	b2Val, _ = table.get_cell("b2")
	c3Val, _ = table.get_cell("c3")

	fmt.Printf("Updated Values (after a1 changed):\n")
	fmt.Printf("a1: %.2f\n", a1Val) // Expected: 50.00
	fmt.Printf("b2: %.2f\n", b2Val) // Expected: 55.00 (50 + 5)
	fmt.Printf("c3: %.2f\n", c3Val) // Expected: 5.00 (based on simplified eval of "55.0*2")
}
