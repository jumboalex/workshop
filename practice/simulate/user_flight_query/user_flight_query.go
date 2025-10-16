package main

import (
	"fmt"
	"sort"
	"time"
)

// Flight 结构体存储单个航班信息
type Flight struct {
	UserID        string
	DepartureCode string
	DepartureTime time.Time
	ArrivalCode   string
	ArrivalTime   time.Time
}

// LocationResult 定义函数的返回结果
type LocationResult struct {
	AirportCode string // 用户所在的机场代码
	Status      string // 用户的状态 (例如: "在机场", "飞行中", "位置未知")
}

// TimeLayout 定义时间字符串的解析格式
const TimeLayout = "2006-01-02 15:04"

// ----------------------------------------------------
// 核心逻辑函数
// ----------------------------------------------------

// FindUserLocation 查找用户在给定时间点的位置
func FindUserLocation(flights []Flight, userID string, checkTime time.Time) LocationResult {
	// 1. 筛选相关航班
	var userFlights []Flight
	for _, f := range flights {
		if f.UserID == userID {
			userFlights = append(userFlights, f)
		}
	}

	if len(userFlights) == 0 {
		return LocationResult{"", "位置未知 (无航班数据)"}
	}

	// 记录最近完成的航班
	var lastFinishedFlight *Flight
	foundFinishedFlight := false

	for i := range userFlights {
		f := &userFlights[i] // 使用指针方便修改和比较

		// 2. 检查用户是否正在飞行
		// 如果检查时间点落在 (出发时间, 到达时间] 区间内
		if !checkTime.Before(f.DepartureTime) && !checkTime.After(f.ArrivalTime) {
			// 在 Go 语言中，我们通常假设到达时间点用户已在到达机场，
			// 但如果在开放区间 (DepartureTime, ArrivalTime) 内，用户就在空中。
			// 这里我们采取更严谨的判断：
			// 如果 checkTime 在 (出发时间, 到达时间) 之间，则在飞行中。
			if checkTime.After(f.DepartureTime) && checkTime.Before(f.ArrivalTime) {
				return LocationResult{
					AirportCode: fmt.Sprintf("%s -> %s", f.DepartureCode, f.ArrivalCode),
					Status:      "飞行中",
				}
			}
			// 如果 checkTime 恰好是 ArrivalTime，我们把它视为已到达
			// 如果 checkTime 恰好是 DepartureTime，我们把它视为在出发机场
		}

		// 3. 查找最近已完成的航班
		// 如果检查时间晚于航班到达时间，说明该航班已完成
		if checkTime.After(f.ArrivalTime) {
			foundFinishedFlight = true
			if lastFinishedFlight == nil || f.ArrivalTime.After(lastFinishedFlight.ArrivalTime) {
				lastFinishedFlight = f
			}
		}

		// 4. 检查是否有即将出发的航班（可能用户就在出发机场）
		// 如果检查时间等于出发时间，用户在出发机场
		if checkTime.Equal(f.DepartureTime) {
			// 如果 checkTime 恰好是某一航班的出发时间，
			// 并且这个航班比最近完成的航班要新，则返回该出发机场
			if !foundFinishedFlight || checkTime.After(lastFinishedFlight.ArrivalTime) {
				return LocationResult{f.DepartureCode, "在机场 (即将出发)"}
			}
		}
	}

	// 5. 返回结果
	if lastFinishedFlight != nil {
		// 返回最近完成航班的到达机场
		return LocationResult{
			AirportCode: lastFinishedFlight.ArrivalCode,
			Status:      "在机场 (已到达)",
		}
	}

	// 如果没有完成的航班，且没有正在飞行的航班，则用户可能在第一个航班的出发机场，
	// 或者根本不在任何一个已知机场。为了严谨，我们返回最早航班的出发机场作为起点。
	sort.Slice(userFlights, func(i, j int) bool {
		return userFlights[i].DepartureTime.Before(userFlights[j].DepartureTime)
	})

	// 如果用户还没开始第一个航班
	if checkTime.Before(userFlights[0].DepartureTime) {
		return LocationResult{userFlights[0].DepartureCode, "在机场 (等待出发)"}
	}

	return LocationResult{"", "位置未知 (数据边界或逻辑异常)"}
}

// ----------------------------------------------------
// 示例运行
// ----------------------------------------------------

func main() {
	// 示例数据 (请确保时间格式与 TimeLayout 一致: YYYY-MM-DD HH:MM)
	allFlights := []Flight{
		{
			UserID: "U123", DepartureCode: "JFK", DepartureTime: mustParseTime("2025-10-09 08:00"),
			ArrivalCode: "ORD", ArrivalTime: mustParseTime("2025-10-09 10:00"),
		},
		{
			UserID: "U123", DepartureCode: "ORD", DepartureTime: mustParseTime("2025-10-09 16:00"),
			ArrivalCode: "LAX", ArrivalTime: mustParseTime("2025-10-09 18:30"),
		},
		{
			UserID: "U123", DepartureCode: "LAX", DepartureTime: mustParseTime("2025-10-10 10:00"),
			ArrivalCode: "MIA", ArrivalTime: mustParseTime("2025-10-10 13:00"),
		},
		{
			UserID: "U123", DepartureCode: "MIA", DepartureTime: mustParseTime("2025-10-10 14:00"),
			ArrivalCode: "DAL", ArrivalTime: mustParseTime("2025-10-10 17:00"),
		},
		{
			UserID: "U456", DepartureCode: "PEK", DepartureTime: mustParseTime("2025-10-10 09:00"),
			ArrivalCode: "PVG", ArrivalTime: mustParseTime("2025-10-10 11:00"),
		},
	}

	// 检查时间点 (2025-10-10 14:30)
	checkTime := mustParseTime("2025-10-10 14:30")
	fmt.Printf("--- 查询时间点: %s ---\n", checkTime.Format(TimeLayout))

	// 案例 1: 正在飞行
	result1 := FindUserLocation(allFlights, "U123", checkTime)
	printResult("U123", result1) // 预期: 飞行中 (MIA -> DAL)

	// 案例 2: 已到达，停留在机场
	checkTime2 := mustParseTime("2025-10-10 13:30")
	result2 := FindUserLocation(allFlights, "U123", checkTime2)
	printResult("U123", result2) // 预期: 在机场 (MIA)

	// 案例 3: 尚未出发
	checkTime3 := mustParseTime("2025-10-08 07:00")
	result3 := FindUserLocation(allFlights, "U123", checkTime3)
	printResult("U123", result3) // 预期: 在机场 (JFK)

	// 案例 4: 无数据
	checkTime4 := mustParseTime("2025-10-10 12:00")
	result4 := FindUserLocation(allFlights, "U789", checkTime4)
	printResult("U789", result4) // 预期: 位置未知

	// 案例 5: 恰好到达 (MIA)
	checkTime5 := mustParseTime("2025-10-10 13:00")
	result5 := FindUserLocation(allFlights, "U123", checkTime5)
	printResult("U123", result5) // 预期: 在机场 (MIA)

	// 案例 6: 恰好出发 (MIA)
	checkTime6 := mustParseTime("2025-10-10 14:00")
	result6 := FindUserLocation(allFlights, "U123", checkTime6)
	printResult("U123", result6) // 预期: 在机场 (MIA)
}

// 辅助函数: 打印结果
func printResult(userID string, res LocationResult) {
	fmt.Printf("用户 %s: 状态: %-16s | 机场/航线: %s\n", userID, res.Status, res.AirportCode)
}

// 辅助函数: 解析时间字符串，如果失败则 panic（仅用于示例）
func mustParseTime(timeStr string) time.Time {
	t, err := time.Parse(TimeLayout, timeStr)
	if err != nil {
		panic(err)
	}
	return t
}
