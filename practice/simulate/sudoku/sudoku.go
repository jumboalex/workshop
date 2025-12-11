package main

func solveSudoku(board [][]byte) {
	backtrack(board, 0, 0)
}

func backtrack(board [][]byte, row, col int) bool {
	// Move to next row if we've finished current row
	if col == 9 {
		return backtrack(board, row+1, 0)
	}

	// All rows filled - solved!
	if row == 9 {
		return true
	}

	// Skip filled cells
	if board[row][col] != '.' {
		return backtrack(board, row, col+1)
	}

	// Try placing numbers
	for n := byte('1'); n <= byte('9'); n++ {
		if canPlace(board, row, col, n) {
			board[row][col] = n
			if backtrack(board, row, col+1) {
				return true
			}
			board[row][col] = '.'
		}
	}

	return false
}

func canPlace(board [][]byte, row, col int, n byte) bool {
	for i := 0; i < 9; i++ {
		if board[i][col] == n {
			return false
		}
		if board[row][i] == n {
			return false
		}
		if board[(row/3)*3+i/3][(col/3)*3+i%3] == n {
			return false
		}
	}
	return true
}
