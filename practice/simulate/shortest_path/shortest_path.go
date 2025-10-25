package main

import (
	"fmt"
	"strings"
)

// Point represents a coordinate in the 2D grid
type Point struct {
	Row int
	Col int
}

// PathNode represents a node in the search with its distance
type PathNode struct {
	Point
	Distance int
}

// ShortestPath finds the shortest distance between two points in a 2D grid
// using bidirectional BFS (searches from both start and end simultaneously)
// The grid contains:
//   0 = walkable cell
//   1 = obstacle/wall
// Returns: distance, found
func ShortestPath(grid [][]int, start, end Point) (int, bool) {
	if !isValid(grid, start) || !isValid(grid, end) {
		return -1, false
	}

	if grid[start.Row][start.Col] == 1 || grid[end.Row][end.Col] == 1 {
		return -1, false
	}

	if start == end {
		return 0, true
	}

	// Directions: up, down, left, right
	directions := []Point{
		{-1, 0}, // up
		{1, 0},  // down
		{0, -1}, // left
		{0, 1},  // right
	}

	// Two queues: one from start, one from end
	queueStart := []*PathNode{{Point: start, Distance: 0}}
	queueEnd := []*PathNode{{Point: end, Distance: 0}}

	// Two visited maps with distances
	visitedStart := make(map[Point]int)
	visitedEnd := make(map[Point]int)
	visitedStart[start] = 0
	visitedEnd[end] = 0

	// Alternate between expanding from start and end
	for len(queueStart) > 0 || len(queueEnd) > 0 {
		// Expand from start
		if len(queueStart) > 0 {
			if dist, found := expandLevel(&queueStart, visitedStart, visitedEnd, grid, directions); found {
				return dist, true
			}
		}

		// Expand from end
		if len(queueEnd) > 0 {
			if dist, found := expandLevel(&queueEnd, visitedEnd, visitedStart, grid, directions); found {
				return dist, true
			}
		}
	}

	// No path found
	return -1, false
}

// expandLevel expands one level of BFS and checks for intersection with the other search
func expandLevel(queue *[]*PathNode, visited, otherVisited map[Point]int, grid [][]int, directions []Point) (int, bool) {
	// Process all nodes at current level
	levelSize := len(*queue)

	for i := 0; i < levelSize; i++ {
		current := (*queue)[0]
		*queue = (*queue)[1:]

		// Explore neighbors
		for _, dir := range directions {
			neighbor := Point{
				Row: current.Row + dir.Row,
				Col: current.Col + dir.Col,
			}

			// Check if neighbor is valid and not an obstacle
			if !isValid(grid, neighbor) || grid[neighbor.Row][neighbor.Col] == 1 {
				continue
			}

			// Check if we've visited this from our side
			if _, seen := visited[neighbor]; seen {
				continue
			}

			newDist := current.Distance + 1
			visited[neighbor] = newDist

			// Check if the other BFS has reached this point
			if otherDist, found := otherVisited[neighbor]; found {
				// Paths meet! Total distance is sum of both distances
				return newDist + otherDist, true
			}

			// Add to queue for further exploration
			*queue = append(*queue, &PathNode{
				Point:    neighbor,
				Distance: newDist,
			})
		}
	}

	return -1, false
}


// isValid checks if a point is within grid bounds
func isValid(grid [][]int, p Point) bool {
	return p.Row >= 0 && p.Row < len(grid) && p.Col >= 0 && p.Col < len(grid[0])
}


// PrintGridWithNumbers prints the grid with coordinates
func PrintGridWithNumbers(grid [][]int) {
	fmt.Print("   ")
	for col := 0; col < len(grid[0]); col++ {
		fmt.Printf("%2d ", col)
	}
	fmt.Println()

	for row := 0; row < len(grid); row++ {
		fmt.Printf("%2d ", row)
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == 1 {
				fmt.Print(" █ ")
			} else {
				fmt.Print(" · ")
			}
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║     2D Grid Shortest Path Finder              ║")
	fmt.Println("╚════════════════════════════════════════════════╝")
	fmt.Println()

	// Test Case 1: Simple grid with obstacles
	grid1 := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0},
		{0, 0, 0, 0, 0},
	}

	start1 := Point{0, 0}
	end1 := Point{4, 4}

	fmt.Println("Test Case 1: Simple Grid")
	fmt.Println("Grid (0 = walkable, 1 = obstacle):")
	PrintGridWithNumbers(grid1)
	fmt.Printf("\nStart: (%d, %d), End: (%d, %d)\n", start1.Row, start1.Col, end1.Row, end1.Col)

	distance, found := ShortestPath(grid1, start1, end1)
	if found {
		fmt.Printf("\n✅ Shortest distance: %d\n", distance)
	} else {
		fmt.Println("\n❌ No path found")
	}

	// Test Case 2: Grid with no path
	fmt.Println("\n" + strings.Repeat("─", 50))
	grid2 := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 1, 0},
	}

	start2 := Point{2, 2}
	end2 := Point{4, 4}

	fmt.Println("\nTest Case 2: Blocked Path")
	fmt.Println("Grid:")
	PrintGridWithNumbers(grid2)
	fmt.Printf("\nStart: (%d, %d), End: (%d, %d)\n", start2.Row, start2.Col, end2.Row, end2.Col)

	distance, found = ShortestPath(grid2, start2, end2)
	if found {
		fmt.Printf("\n✅ Shortest distance: %d\n", distance)
	} else {
		fmt.Println("\n❌ No path found")
	}

	// Test Case 3: Multiple paths
	fmt.Println("\n" + strings.Repeat("─", 50))
	grid3 := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
	}

	start3 := Point{0, 0}
	end3 := Point{4, 4}

	fmt.Println("\nTest Case 3: Grid with Multiple Possible Paths")
	fmt.Println("Grid:")
	PrintGridWithNumbers(grid3)
	fmt.Printf("\nStart: (%d, %d), End: (%d, %d)\n", start3.Row, start3.Col, end3.Row, end3.Col)

	distance, found = ShortestPath(grid3, start3, end3)
	if found {
		fmt.Printf("\n✅ Shortest distance: %d\n", distance)
	}

	// Test Case 4: Large maze
	fmt.Println("\n" + strings.Repeat("─", 50))
	grid4 := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 0, 1, 1, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	start4 := Point{0, 0}
	end4 := Point{6, 9}

	fmt.Println("\nTest Case 4: Larger Maze")
	fmt.Println("Grid:")
	PrintGridWithNumbers(grid4)
	fmt.Printf("\nStart: (%d, %d), End: (%d, %d)\n", start4.Row, start4.Col, end4.Row, end4.Col)

	distance, found = ShortestPath(grid4, start4, end4)
	if found {
		fmt.Printf("\n✅ Shortest distance: %d\n", distance)
	} else {
		fmt.Println("\n❌ No path found")
	}

	// Test Case 5: Same start and end
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("\nTest Case 5: Start = End")
	start5 := Point{2, 2}
	end5 := Point{2, 2}
	distance, found = ShortestPath(grid1, start5, end5)
	if found {
		fmt.Printf("✅ Distance: %d (same position)\n", distance)
	}
}
