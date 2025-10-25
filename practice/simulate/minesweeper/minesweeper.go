package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Cell represents a single cell in the minesweeper grid
type Cell struct {
	IsMine     bool
	IsRevealed bool
	IsFlagged  bool
	AdjacentMines int
}

// Game represents the minesweeper game state
type Game struct {
	Width      int
	Height     int
	MineCount  int
	Grid       [][]Cell
	GameOver   bool
	Won        bool
	StartTime  time.Time
}

// NewGame creates a new minesweeper game
func NewGame(width, height, mineCount int) *Game {
	game := &Game{
		Width:     width,
		Height:    height,
		MineCount: mineCount,
		Grid:      make([][]Cell, height),
	}

	// Initialize grid
	for i := range game.Grid {
		game.Grid[i] = make([]Cell, width)
	}

	return game
}

// PlaceMines randomly places mines on the board, avoiding the first clicked cell
func (g *Game) PlaceMines(firstRow, firstCol int) {
	rand.Seed(time.Now().UnixNano())

	placed := 0
	for placed < g.MineCount {
		row := rand.Intn(g.Height)
		col := rand.Intn(g.Width)

		// Don't place mine on first clicked cell or if already has mine
		if (row == firstRow && col == firstCol) || g.Grid[row][col].IsMine {
			continue
		}

		g.Grid[row][col].IsMine = true
		placed++
	}

	// Calculate adjacent mine counts
	g.calculateAdjacentMines()
}

// calculateAdjacentMines counts mines adjacent to each cell
func (g *Game) calculateAdjacentMines() {
	for row := 0; row < g.Height; row++ {
		for col := 0; col < g.Width; col++ {
			if !g.Grid[row][col].IsMine {
				g.Grid[row][col].AdjacentMines = g.countAdjacentMines(row, col)
			}
		}
	}
}

// countAdjacentMines counts how many mines are adjacent to a cell
func (g *Game) countAdjacentMines(row, col int) int {
	count := 0
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1},           {0, 1},
		{1, -1},  {1, 0},  {1, 1},
	}

	for _, dir := range directions {
		newRow := row + dir[0]
		newCol := col + dir[1]

		if g.isValidCell(newRow, newCol) && g.Grid[newRow][newCol].IsMine {
			count++
		}
	}

	return count
}

// isValidCell checks if coordinates are within the grid
func (g *Game) isValidCell(row, col int) bool {
	return row >= 0 && row < g.Height && col >= 0 && col < g.Width
}

// RevealCell reveals a cell and potentially cascades to adjacent cells
func (g *Game) RevealCell(row, col int) {
	if !g.isValidCell(row, col) || g.Grid[row][col].IsRevealed || g.Grid[row][col].IsFlagged {
		return
	}

	cell := &g.Grid[row][col]
	cell.IsRevealed = true

	// Hit a mine - game over
	if cell.IsMine {
		g.GameOver = true
		g.Won = false
		return
	}

	// If no adjacent mines, reveal neighbors (flood fill)
	if cell.AdjacentMines == 0 {
		directions := [][]int{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1},           {0, 1},
			{1, -1},  {1, 0},  {1, 1},
		}

		for _, dir := range directions {
			g.RevealCell(row+dir[0], col+dir[1])
		}
	}
}

// ToggleFlag toggles a flag on a cell
func (g *Game) ToggleFlag(row, col int) {
	if !g.isValidCell(row, col) || g.Grid[row][col].IsRevealed {
		return
	}

	g.Grid[row][col].IsFlagged = !g.Grid[row][col].IsFlagged
}

// CheckWin checks if the player has won
func (g *Game) CheckWin() bool {
	for row := 0; row < g.Height; row++ {
		for col := 0; col < g.Width; col++ {
			cell := g.Grid[row][col]
			// If there's a non-mine cell that's not revealed, game is not won
			if !cell.IsMine && !cell.IsRevealed {
				return false
			}
		}
	}
	return true
}

// Display prints the current game board
func (g *Game) Display(showMines bool) {
	fmt.Println()

	// Print column numbers
	fmt.Print("    ")
	for col := 0; col < g.Width; col++ {
		fmt.Printf("%2d ", col)
	}
	fmt.Println()

	// Print top border
	fmt.Print("   ╔")
	for col := 0; col < g.Width; col++ {
		fmt.Print("══")
		if col < g.Width-1 {
			fmt.Print("═")
		}
	}
	fmt.Println("═╗")

	// Print grid
	for row := 0; row < g.Height; row++ {
		fmt.Printf("%2d ║ ", row)

		for col := 0; col < g.Width; col++ {
			cell := g.Grid[row][col]

			if cell.IsFlagged && !showMines {
				fmt.Print("🚩")
			} else if !cell.IsRevealed && !showMines {
				fmt.Print("▢ ")
			} else if cell.IsMine {
				if showMines {
					fmt.Print("💣")
				} else {
					fmt.Print("💣")
				}
			} else if cell.AdjacentMines == 0 {
				fmt.Print("  ")
			} else {
				// Color code the numbers
				fmt.Printf("%s%d%s ", getNumberColor(cell.AdjacentMines), cell.AdjacentMines, resetColor())
			}

			if col < g.Width-1 {
				fmt.Print(" ")
			}
		}

		fmt.Println("║")
	}

	// Print bottom border
	fmt.Print("   ╚")
	for col := 0; col < g.Width; col++ {
		fmt.Print("══")
		if col < g.Width-1 {
			fmt.Print("═")
		}
	}
	fmt.Println("═╝")
	fmt.Println()
}

// getNumberColor returns ANSI color code for mine count numbers
func getNumberColor(num int) string {
	colors := map[int]string{
		1: "\033[34m", // Blue
		2: "\033[32m", // Green
		3: "\033[31m", // Red
		4: "\033[35m", // Purple
		5: "\033[33m", // Yellow
		6: "\033[36m", // Cyan
		7: "\033[90m", // Gray
		8: "\033[91m", // Bright Red
	}
	if color, ok := colors[num]; ok {
		return color
	}
	return ""
}

// resetColor returns ANSI reset code
func resetColor() string {
	return "\033[0m"
}

// GetFlagCount returns the number of flags placed
func (g *Game) GetFlagCount() int {
	count := 0
	for row := 0; row < g.Height; row++ {
		for col := 0; col < g.Width; col++ {
			if g.Grid[row][col].IsFlagged {
				count++
			}
		}
	}
	return count
}

// main game loop
func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║      MINESWEEPER GAME 💣           ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()

	// Select difficulty
	fmt.Println("Select difficulty:")
	fmt.Println("1. Beginner    (9x9,  10 mines)")
	fmt.Println("2. Intermediate (16x16, 40 mines)")
	fmt.Println("3. Expert      (30x16, 99 mines)")
	fmt.Println("4. Custom")
	fmt.Print("\nChoice (1-4): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var width, height, mines int

	switch choice {
	case "1":
		width, height, mines = 9, 9, 10
	case "2":
		width, height, mines = 16, 16, 40
	case "3":
		width, height, mines = 30, 16, 99
	case "4":
		fmt.Print("Width: ")
		w, _ := reader.ReadString('\n')
		width, _ = strconv.Atoi(strings.TrimSpace(w))

		fmt.Print("Height: ")
		h, _ := reader.ReadString('\n')
		height, _ = strconv.Atoi(strings.TrimSpace(h))

		fmt.Print("Mines: ")
		m, _ := reader.ReadString('\n')
		mines, _ = strconv.Atoi(strings.TrimSpace(m))
	default:
		width, height, mines = 9, 9, 10
	}

	game := NewGame(width, height, mines)
	firstMove := true

	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║           HOW TO PLAY              ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ Commands:                          ║")
	fmt.Println("║   r <row> <col>  - Reveal cell     ║")
	fmt.Println("║   f <row> <col>  - Flag/unflag     ║")
	fmt.Println("║   q              - Quit            ║")
	fmt.Println("╚════════════════════════════════════╝")

	for !game.GameOver {
		game.Display(false)

		flagCount := game.GetFlagCount()
		remainingMines := game.MineCount - flagCount
		fmt.Printf("Mines remaining: %d | Flags: %d/%d\n", remainingMines, flagCount, game.MineCount)

		if firstMove {
			fmt.Println("\n💡 First move is always safe!")
		}

		fmt.Print("\nCommand: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]

		if cmd == "q" {
			fmt.Println("\nThanks for playing!")
			return
		}

		if len(parts) < 3 {
			fmt.Println("❌ Invalid command. Use: r <row> <col> or f <row> <col>")
			continue
		}

		row, err1 := strconv.Atoi(parts[1])
		col, err2 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil {
			fmt.Println("❌ Invalid coordinates")
			continue
		}

		if !game.isValidCell(row, col) {
			fmt.Println("❌ Coordinates out of bounds")
			continue
		}

		switch cmd {
		case "r":
			if firstMove {
				game.PlaceMines(row, col)
				game.StartTime = time.Now()
				firstMove = false
			}
			game.RevealCell(row, col)

			if game.GameOver {
				game.Display(true)
				fmt.Println("💥 BOOM! You hit a mine!")
				fmt.Println("╔════════════════════════════════════╗")
				fmt.Println("║          GAME OVER 💀              ║")
				fmt.Println("╚════════════════════════════════════╝")
			} else if game.CheckWin() {
				game.GameOver = true
				game.Won = true
				elapsed := time.Since(game.StartTime)
				game.Display(true)
				fmt.Println("🎉 Congratulations! You won!")
				fmt.Printf("⏱️  Time: %.1f seconds\n", elapsed.Seconds())
				fmt.Println("╔════════════════════════════════════╗")
				fmt.Println("║          VICTORY! 🏆               ║")
				fmt.Println("╚════════════════════════════════════╝")
			}

		case "f":
			game.ToggleFlag(row, col)

		default:
			fmt.Println("❌ Invalid command. Use 'r' to reveal or 'f' to flag")
		}
	}
}
