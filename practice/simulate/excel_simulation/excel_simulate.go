package main

import (
	"fmt"
	"regexp"
	"sort"
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
	if oldCell, exists := t.cells[cell]; exists {
		for _, dep := range oldCell.deps {
			t.removeRevDep(dep, cell)
		}
	}

	// 2. Parse and validate formula
	cellStruct, err := t.parseFormula(formula)
	if err != nil {
		return fmt.Errorf("cell %s: %w", cell, err)
	}

	// 3. Update reverse dependencies
	for _, dep := range cellStruct.deps {
		t.addRevDep(dep, cell)
	}

	// 4. Store and recalculate
	t.cells[cell] = cellStruct
	t.recalculate(cell)

	return nil
}

// parseFormula parses a formula string and returns a Cell struct
func (t *Table) parseFormula(formula string) (*Cell, error) {
	cell := &Cell{formula: formula}

	// Try parsing as a number
	if val, err := strconv.ParseFloat(formula, 64); err == nil {
		cell.value = val
		return cell, nil
	}

	// Try parsing as Excel-style formula (with "=")
	if formulaContent, hasPrefix := strings.CutPrefix(formula, "="); hasPrefix {
		cell.formula = formulaContent
		cell.deps = cellRefRegex.FindAllString(formulaContent, -1)
		return cell, nil
	}

	// Invalid format
	return nil, fmt.Errorf("invalid formula format: formulas must start with '=' (e.g., '=A1+B1')")
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
	cellStruct, exists := t.cells[cell]
	if !exists {
		return
	}

	// If cell has dependencies, evaluate the formula
	if len(cellStruct.deps) > 0 {
		// Build value map for replacement
		valueMap := make(map[string]float64)
		for _, dep := range cellStruct.deps {
			depCell, exists := t.cells[dep]
			if exists {
				valueMap[dep] = depCell.value
			} else {
				valueMap[dep] = 0.0
			}
		}

		// Replace cell references with values (sorted by length descending to avoid substring issues)
		evalFormula := t.replaceCellRefsWithValues(cellStruct.formula, valueMap)

		// Evaluate the expression
		result := t.evaluateExpression(evalFormula)
		cellStruct.value = result
	}

	// Propagate changes to dependent cells
	for _, dependentCell := range t.revDeps[cell] {
		t.recalculate(dependentCell)
	}
}

// replaceCellRefsWithValues replaces cell references with their values
// Sorts cell names by length (descending) to avoid substring conflicts (e.g., a10 before a1)
func (t *Table) replaceCellRefsWithValues(formula string, valueMap map[string]float64) string {
	// Sort cell names by length (longest first) to avoid substring issues
	cellNames := make([]string, 0, len(valueMap))
	for name := range valueMap {
		cellNames = append(cellNames, name)
	}

	// Sort by length descending
	sort.Slice(cellNames, func(i, j int) bool {
		return len(cellNames[i]) > len(cellNames[j])
	})

	result := formula
	for _, name := range cellNames {
		value := valueMap[name]
		result = strings.ReplaceAll(result, name, strconv.FormatFloat(value, 'f', -1, 64))
	}
	return result
}

// evaluateExpression evaluates a simple math expression using a stack-based approach
// Supports +, -, *, / with proper precedence
func (t *Table) evaluateExpression(expression string) float64 {
	expression = strings.TrimSpace(expression)

	// Try parsing as simple number first
	val, err := strconv.ParseFloat(expression, 64)
	if err == nil {
		return val
	}

	tokens := t.tokenize(expression)
	if len(tokens) == 0 {
		return 0
	}

	// Two stacks: one for values, one for operators
	values := []float64{}
	ops := []string{}

	precedence := func(op string) int {
		if op == "+" || op == "-" {
			return 1
		}
		if op == "*" || op == "/" {
			return 2
		}
		return 0
	}

	applyOp := func(a, b float64, op string) float64 {
		switch op {
		case "+":
			return a + b
		case "-":
			return a - b
		case "*":
			return a * b
		case "/":
			if b != 0 {
				return a / b
			}
			return 0
		}
		return 0
	}

	for _, token := range tokens {
		// If token is a number, push to values stack
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			values = append(values, num)
		} else {
			// Token is an operator
			// While top of ops has same or greater precedence, apply it
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(token) {
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				if len(values) >= 2 {
					b := values[len(values)-1]
					a := values[len(values)-2]
					values = values[:len(values)-2]
					values = append(values, applyOp(a, b, op))
				}
			}
			ops = append(ops, token)
		}
	}

	// Apply remaining operators
	for len(ops) > 0 {
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		if len(values) >= 2 {
			b := values[len(values)-1]
			a := values[len(values)-2]
			values = values[:len(values)-2]
			values = append(values, applyOp(a, b, op))
		}
	}

	if len(values) > 0 {
		return values[0]
	}
	return 0
}

// tokenize splits expression into numbers and operators
func (t *Table) tokenize(expr string) []string {
	tokens := []string{}
	current := ""

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(ch))
		} else if ch != ' ' {
			current += string(ch)
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
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

	// 2. Set a dependent formula (Excel-style with "=")
	table.set_cell("b2", "=a1+5") // b2 = 10 + 5 = 15
	// Dependencies: b2 depends on a1

	// 3. Set a second-level dependent formula (Excel-style with "=")
	table.set_cell("c3", "=b2*2") // c3 = 15 * 2 = 30
	// Dependencies: c3 depends on b2, which depends on a1

	a1Val, _ := table.get_cell("a1")
	b2Val, _ := table.get_cell("b2")
	c3Val, _ := table.get_cell("c3")

	fmt.Printf("Initial Values:\n")
	fmt.Printf("a1: %.2f\n", a1Val) // Expected: 10.00
	fmt.Printf("b2: %.2f\n", b2Val) // Expected: 15.00 (10 + 5)
	fmt.Printf("c3: %.2f\n", c3Val) // Expected: 30.00 (15 * 2)
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
	fmt.Printf("c3: %.2f\n", c3Val) // Expected: 110.00 (55 * 2)
	fmt.Println("--------------------")

	// 5. Complex formula with multiple cell references
	fmt.Printf("Setting d4 to '=a1+b2'...\n")
	table.set_cell("d4", "=a1+b2")
	// Dependencies: d4 depends on both a1 and b2

	d4Val, _ := table.get_cell("d4")
	fmt.Printf("Complex Formula Test:\n")
	fmt.Printf("d4 (a1+b2): %.2f\n", d4Val) // Expected: 105.00 (50 + 55)
	fmt.Println("--------------------")

	// 6. Update a1 again to verify d4 recalculates
	fmt.Printf("Updating a1 to '100'...\n")
	table.set_cell("a1", "100")

	a1Val, _ = table.get_cell("a1")
	b2Val, _ = table.get_cell("b2")
	c3Val, _ = table.get_cell("c3")
	d4Val, _ = table.get_cell("d4")

	fmt.Printf("Final Values (after a1 changed to 100):\n")
	fmt.Printf("a1: %.2f\n", a1Val)  // Expected: 100.00
	fmt.Printf("b2: %.2f\n", b2Val)  // Expected: 105.00 (100 + 5)
	fmt.Printf("c3: %.2f\n", c3Val)  // Expected: 210.00 (105 * 2)
	fmt.Printf("d4: %.2f\n", d4Val)  // Expected: 205.00 (100 + 105)
	fmt.Println("--------------------")

	// 7. Test error handling for formulas without "="
	fmt.Printf("Testing error handling: Setting e5 to 'a1*3' (without = prefix)...\n")
	err := table.set_cell("e5", "a1*3")
	if err != nil {
		fmt.Printf("Error (as expected): %v\n", err)
	}
	fmt.Println("--------------------")

	// 8. Test correct Excel-style formula
	fmt.Printf("Setting e5 to '=a1*3' (with = prefix)...\n")
	table.set_cell("e5", "=a1*3")
	e5Val, _ := table.get_cell("e5")
	fmt.Printf("Excel-style Formula Test:\n")
	fmt.Printf("e5 (=a1*3): %.2f\n", e5Val) // Expected: 300.00 (100 * 3)
}
