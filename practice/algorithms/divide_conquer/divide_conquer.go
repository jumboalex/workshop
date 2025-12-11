package divideconquer

func searchMatrix(matrix [][]int, target int) bool {
	r := len(matrix)
	c := len(matrix[0])

	return search(matrix, 0, 0, r-1, c-1, target)
}

func search(matrix [][]int, row, col, endR, endC, target int) bool {
	if row > endR || col > endC {
		return false
	}

	midR := (row + endR) / 2
	midC := (col + endC) / 2

	if matrix[midR][midC] == target {
		return true
	}

	if matrix[midR][midC] < target {
		return search(matrix, row, midC+1, midR, endC, target) || search(matrix, midR+1, col, endR, endC, target)
	}
	return search(matrix, row, col, endR, midC-1, target) || search(matrix, row, midC, midR-1, endC, target)
}
