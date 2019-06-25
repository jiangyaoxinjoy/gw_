package utils

import (
	"time"
)

func TimeToTimestamp(toBeCharge string) (int64, error) {
	timeLayout := "2006-01-02 15:04:05"

	// toBeCharge := timeMap[0]["startTime"]
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	startTime := theTime.Unix()                                       //转化为时间戳 类型是int64
	if err != nil {
		return startTime, err
	}
	return startTime, nil
	// fmt.Println(theTime)                                            //打印输出theTime 2015-01-01 15:15:00 +0800 CST
	// fmt.Println(sr)                                                 //打印输出时间戳 1420041600
	// fmt.Println(timeMap[0]["startTime"])
}
