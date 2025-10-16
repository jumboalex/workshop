package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ----------------------------------------------------
// 1. 数据结构定义 (用于解析 API 响应)
// ----------------------------------------------------

// ExchangeRateAPIResponse 定义了 ExchangeRate.host 历史汇率 API 的 JSON 响应结构
type ExchangeRateAPIResponse struct {
	Success bool               `json:"success"`
	Date    string             `json:"date"`
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
	Error   struct {
		Code int    `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

// ----------------------------------------------------
// 2. 命令行工具逻辑
// ----------------------------------------------------

func main() {
	// 定义命令行参数
	startDateStr := flag.String("start-date", "", "查询的开始日期 (格式: YYYY-MM-DD)")
	endDateStr := flag.String("end-date", "", "查询的结束日期 (格式: YYYY-MM-DD)")
	fromCurrency := flag.String("from", "USD", "基准货币 (例如: USD, EUR)")
	toCurrency := flag.String("to", "CNY", "目标货币 (例如: CNY, JPY)")

	flag.Parse()

	// 基础验证
	if *startDateStr == "" || *endDateStr == "" {
		fmt.Println("❌ 错误: 必须指定开始日期和结束日期。")
		flag.Usage()
		os.Exit(1)
	}

	// 解析日期
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, *startDateStr)
	if err != nil {
		fmt.Printf("❌ 错误: 无效的开始日期格式。请使用 YYYY-MM-DD。 %v\n", err)
		os.Exit(1)
	}
	endDate, err := time.Parse(layout, *endDateStr)
	if err != nil {
		fmt.Printf("❌ 错误: 无效的结束日期格式。请使用 YYYY-MM-DD。 %v\n", err)
		os.Exit(1)
	}

	if startDate.After(endDate) {
		fmt.Println("❌ 错误: 开始日期不能晚于结束日期。")
		os.Exit(1)
	}

	fmt.Printf("📊 查询汇率: %s -> %s (从 %s 到 %s)\n",
		strings.ToUpper(*fromCurrency),
		strings.ToUpper(*toCurrency),
		startDate.Format(layout),
		endDate.Format(layout),
	)
	fmt.Println("-------------------------------------------------")
	fmt.Printf("%-12s | %s 兑 1 %s\n", "日期", strings.ToUpper(*toCurrency), strings.ToUpper(*fromCurrency))
	fmt.Println("-------------------------------------------------")

	// 迭代日期并获取汇率
	currentDate := startDate
	for !currentDate.After(endDate) {
		rate, err := getHistoricalRate(currentDate, *fromCurrency, *toCurrency)
		if err != nil {
			fmt.Printf("%s | ❌ 获取失败: %v\n", currentDate.Format(layout), err)
		} else {
			fmt.Printf("%s | %.4f\n", currentDate.Format(layout), rate)
		}

		// 移动到下一天
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	fmt.Println("-------------------------------------------------")
}

// ----------------------------------------------------
// 3. API 调用函数
// ----------------------------------------------------

// getHistoricalRate 调用 API 获取特定日期、基准货币和目标货币的汇率
func getHistoricalRate(date time.Time, base string, symbol string) (float64, error) {
	// API URL 格式: https://api.exchangerate.host/YYYY-MM-DD?base=BASE&symbols=SYMBOLS
	// 注意: ExchangeRate.host 免费层级通常不需要 API Key
	dateStr := date.Format("2006-01-02")
	apiURL := fmt.Sprintf("https://api.exchangerate.host/%s?base=%s&symbols=%s",
		dateStr,
		strings.ToUpper(base),
		strings.ToUpper(symbol))

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API 返回状态码 %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ExchangeRateAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("JSON 解析失败: %w", err)
	}

	// 检查 API 业务逻辑错误
	if !result.Success {
		return 0, fmt.Errorf("API 错误代码 %d: %s", result.Error.Code, result.Error.Info)
	}

	// 提取汇率
	rate, ok := result.Rates[strings.ToUpper(symbol)]
	if !ok {
		// 如果目标货币不在返回的 Rates map 中，通常表示货币代码错误
		return 0, fmt.Errorf("在 API 响应中找不到目标货币 %s 的汇率", strings.ToUpper(symbol))
	}

	return rate, nil
}
