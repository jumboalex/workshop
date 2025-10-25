package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

// ----------------------------------------------------
// Step 1: 读取本地词库
// ----------------------------------------------------

// loadVocabulary 从指定路径加载词库文件，返回一个包含所有单词的 map（set）。
func loadVocabulary(filePath string) (map[string]struct{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开词库文件失败: %w", err)
	}
	defer file.Close()

	vocabulary := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if word != "" {
			vocabulary[word] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取词库文件内容失败: %w", err)
	}
	fmt.Printf("✅ 成功读取词库，包含 %d 个单词。\n", len(vocabulary))
	return vocabulary, nil
}

// ----------------------------------------------------
// Step 2: 从 URL 下载并读取文本
// ----------------------------------------------------

// downloadText 从 URL 下载文本并返回内容。
func downloadText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 请求返回状态码 %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %w", err)
	}

	rawText := strings.ToLower(string(bodyBytes))
	fmt.Printf("✅ 成功下载文本文件，内容大小: %d 字符。\n", len(rawText))
	return rawText, nil
}

// ----------------------------------------------------
// Step 3: 清理、过滤单词并计算长度分布
// ----------------------------------------------------

// calculateDistribution 计算单词长度分布。
func calculateDistribution(text string, vocabulary map[string]struct{}) map[int]int {
	// 匹配所有连续的字母字符 (a-z)
	re := regexp.MustCompile("[^a-z]+")
	cleanedText := re.ReplaceAllString(text, " ")

	potentialWords := strings.Fields(cleanedText)

	lengthCounts := make(map[int]int)

	for _, word := range potentialWords {
		// 检查单词是否在词库中
		if _, exists := vocabulary[word]; exists {
			length := len(word)
			lengthCounts[length]++
		}
	}

	return lengthCounts
}

// ----------------------------------------------------
// Step 4: 打印直方图
// ----------------------------------------------------

// printHistogram 打印单词长度分布的直方图。
func printHistogram(lengthCounts map[int]int) {
	if len(lengthCounts) == 0 {
		fmt.Println("\n⚠️ 警告: 文本中没有找到任何在词库中的单词。无法生成分布。")
		return
	}

	var totalValidWords int
	var maxCount int
	var sortedLengths []int

	// 计算总词数，最大计数，并收集所有长度
	for length, count := range lengthCounts {
		totalValidWords += count
		if count > maxCount {
			maxCount = count
		}
		sortedLengths = append(sortedLengths, length)
	}

	// 排序长度 (从小到大)
	sort.Ints(sortedLengths)

	fmt.Printf("\n📊 总共找到 %d 个符合词库条件的单词。\n", totalValidWords)
	fmt.Println("--- 单词长度分布 ---")

	barWidth := 50 // 直方图的最大显示宽度

	for _, length := range sortedLengths {
		count := lengthCounts[length]
		// 计算百分比
		percentage := float64(count) / float64(totalValidWords) * 100
		// 计算直方图的块数
		stars := int((float64(count) / float64(maxCount)) * float64(barWidth))

		// 打印结果
		fmt.Printf("  %2d | %5d (%5.1f%%) | %s\n",
			length,
			count,
			percentage,
			strings.Repeat("█", stars))
	}

	fmt.Println("--------------------")
	fmt.Printf("（直方图最大宽度: %d，█ 代表比例）\n", barWidth)
}

// ----------------------------------------------------
// 主函数
// ----------------------------------------------------

func main() {
	// 示例 URL (替换成您自己的文本文件链接)
	textURL := "https://www.gutenberg.org/files/11/11-0.txt" // 爱丽丝梦游仙境

	// 示例词库文件路径 (请确保在本地创建了此文件，例如 local_vocab.txt)
	// 文件内容示例:
	// the
	// is
	// and
	// rabbit
	// alice
	vocabularyFilePath := "local_vocab.txt"

	// 1. 加载词库
	vocabulary, err := loadVocabulary(vocabularyFilePath)
	if err != nil {
		fmt.Printf("程序终止: %v\n", err)
		return
	}

	// 2. 下载文本
	text, err := downloadText(textURL)
	if err != nil {
		fmt.Printf("程序终止: %v\n", err)
		return
	}

	// 3. 计算分布
	lengthCounts := calculateDistribution(text, vocabulary)

	// 4. 打印直方图
	printHistogram(lengthCounts)
}
