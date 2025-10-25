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

	// 2. 按出发时间排序航班
	sort.Slice(userFlights, func(i, j int) bool {
		return userFlights[i].DepartureTime.Before(userFlights[j].DepartureTime)
	})

	// 3. 使用二分查找找到最后一个出发时间 <= checkTime 的航班
	// 这样可以快速定位用户可能所在的航班或已完成的最近航班
	idx := sort.Search(len(userFlights), func(i int) bool {
		return userFlights[i].DepartureTime.After(checkTime)
	})
	// idx 现在指向第一个出发时间 > checkTime 的航班
	// 所以 idx-1 是最后一个出发时间 <= checkTime 的航班

	// 4. 边界情况：如果所有航班都还未出发
	if idx == 0 {
		return LocationResult{userFlights[0].DepartureCode, "在机场 (等待出发)"}
	}

	// 5. 检查最近的已出发航班（userFlights[idx-1]）
	recentFlight := &userFlights[idx-1]

	// 5a. 检查用户是否正在飞行中
	if checkTime.After(recentFlight.DepartureTime) && checkTime.Before(recentFlight.ArrivalTime) {
		return LocationResult{
			AirportCode: fmt.Sprintf("%s -> %s", recentFlight.DepartureCode, recentFlight.ArrivalCode),
			Status:      "飞行中",
		}
	}

	// 5b. 检查是否恰好在出发时刻
	if checkTime.Equal(recentFlight.DepartureTime) {
		return LocationResult{recentFlight.DepartureCode, "在机场 (即将出发)"}
	}

	// 5c. 检查是否恰好在到达时刻（视为已到达）
	if checkTime.Equal(recentFlight.ArrivalTime) {
		return LocationResult{recentFlight.ArrivalCode, "在机场 (已到达)"}
	}

	// 6. 如果 checkTime 在 recentFlight 的到达时间之后
	// 说明用户已完成该航班，停留在到达机场
	if checkTime.After(recentFlight.ArrivalTime) {
		return LocationResult{
			AirportCode: recentFlight.ArrivalCode,
			Status:      "在机场 (已到达)",
		}
	}

	// 7. 理论上不应该到达这里
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

	fmt.Println("\n--- 额外测试边界情况 ---")

	// 案例 7: 在两个航班之间的间隔时间（停留在机场）
	checkTime7 := mustParseTime("2025-10-09 12:00")
	result7 := FindUserLocation(allFlights, "U123", checkTime7)
	printResult("U123", result7) // 预期: 在机场 (ORD, 已到达)

	// 案例 8: 正在第一个航班飞行中
	checkTime8 := mustParseTime("2025-10-09 09:00")
	result8 := FindUserLocation(allFlights, "U123", checkTime8)
	printResult("U123", result8) // 预期: 飞行中 (JFK -> ORD)

	// 案例 9: 完成所有航班后
	checkTime9 := mustParseTime("2025-10-11 00:00")
	result9 := FindUserLocation(allFlights, "U123", checkTime9)
	printResult("U123", result9) // 预期: 在机场 (DAL, 已到达)
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
