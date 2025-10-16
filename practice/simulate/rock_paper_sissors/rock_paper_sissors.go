package main

import (
	"errors"
	"fmt"
	"strconv"
)

// 游戏规则的常量
const (
	Rock     = 'R'
	Paper    = 'P'
	Scissors = 'S'
)

// DetermineWinner 解析剪刀石头布游戏结果字符串并返回赢家。
// 输入格式：[所需胜场数][局数1P1出招][局数1P2出招][局数2P1出招][局数2P2出招]...
func DetermineWinner(resultStr string) (string, error) {
	if len(resultStr) < 1 {
		return "", errors.New("无效输入：字符串为空")
	}

	// 1. 解析所需胜场数 (Best-of-N)
	// 第一个字符必须是数字
	nStr := string(resultStr[0])
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		return "", errors.New("无效输入：第一个字符必须是大于0的数字，表示所需胜场数")
	}

	// 剩余的字符串是游戏结果
	gameResults := resultStr[1:]
	if len(gameResults)%2 != 0 {
		return "", errors.New("无效输入：游戏结果必须成对出现（P1出招, P2出招）")
	}

	player1Score := 0
	player2Score := 0

	// 2. 遍历游戏结果并计算分数
	for i := 0; i < len(gameResults); i += 2 {
		p1Move := gameResults[i]
		p2Move := gameResults[i+1]

		// 验证出招是否有效 (R, P, S)
		if !isValidMove(p1Move) || !isValidMove(p2Move) {
			return "", fmt.Errorf("无效输入：在第 %d 局发现无效出招 ('%c' 或 '%c')，有效出招仅为 R, P, S", (i/2)+1, p1Move, p2Move)
		}

		// 判断胜负
		winner := checkRoundWinner(p1Move, p2Move)

		if winner == 1 {
			player1Score++
		} else if winner == 2 {
			player2Score++
		}
		// 如果是平局 (winner == 0)，则分数不变

		// 提前检查：如果有玩家已经达到所需胜场数，则可立即宣布胜利
		if player1Score >= n {
			return "Player 1 赢了！", nil
		}
		if player2Score >= n {
			return "Player 2 赢了！", nil
		}
	}

	// 3. 判断最终结果
	// 比赛结束但无人达到所需胜场数
	if player1Score > player2Score {
		return "平局：Player 1 领先，但双方均未达到所需胜场数 (" + strconv.Itoa(n) + ")", nil
	} else if player2Score > player1Score {
		return "平局：Player 2 领先，但双方均未达到所需胜场数 (" + strconv.Itoa(n) + ")", nil
	} else {
		return "平局：双方得分相同，且均未达到所需胜场数 (" + strconv.Itoa(n) + ")", nil
	}
}

// isValidMove 检查字符是否是有效的出招。
func isValidMove(move byte) bool {
	return move == Rock || move == Paper || move == Scissors
}

// checkRoundWinner 判断一局的胜者。
// 返回 1 表示 Player 1 赢， 2 表示 Player 2 赢， 0 表示平局。
func checkRoundWinner(p1Move, p2Move byte) int {
	if p1Move == p2Move {
		return 0 // 平局
	}

	switch p1Move {
	case Rock:
		if p2Move == Scissors {
			return 1 // 石头胜剪刀
		}
	case Paper:
		if p2Move == Rock {
			return 1 // 布胜石头
		}
	case Scissors:
		if p2Move == Paper {
			return 1 // 剪刀胜布
		}
	}

	return 2 // 否则 Player 2 赢
}

// --- 示例用法 ---

func main() {
	// 示例：5局3胜，R: 石头, P: 布, S: 剪刀
	examples := map[string]string{
		"5RPSS":      "示例1：5局3胜, P1: R, P2: P; P1: S, P2: S。",
		"3RPPRR":     "示例2：3局2胜, P1: R, P2: P; P1: P, P2: R; P1: R, P2: R。",               // P1 (2) - P2 (1) -> P1 Win
		"5RPSSPRP":   "示例3：5局3胜, P1: R, P2: P; P1: S, P2: S; P1: P, P2: R; P1: P, P2: R。", // P1 (2) - P2 (1) -> 平局 (未达3胜)
		"3RP":        "示例4：3局2胜，P1: R, P2: P。 (未完成)",
		"1PR":        "示例5：1局1胜，P1: P, P2: R。",
		"3RPA":       "示例6：无效出招 'A'。",
		"0RPRP":      "示例7：无效胜场数 '0'。",
		"5RPRPSRPRP": "示例8：5局3胜，P1先达到3胜。", // P1 (3) - P2 (2) -> P1 Win
	}

	for input, desc := range examples {
		fmt.Println("\n----------------------------------------------------")
		fmt.Printf("%s\n输入: '%s'\n", desc, input)
		winner, err := DetermineWinner(input)

		if err != nil {
			fmt.Printf("结果: 错误 -> %v\n", err)
		} else {
			fmt.Printf("结果: %s\n", winner)
		}
	}
}
