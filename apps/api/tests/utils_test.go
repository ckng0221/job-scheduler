package tests

import (
	"fmt"
	"job-scheduler/utils"
	"testing"
	"time"
)

func TestGetUnixMinuteRange(t *testing.T) {
	testCases := []struct {
		input    int64
		minStart int64
		minEnd   int64
	}{

		{1710900325, 1710900300, 1710900360},
		{1704038430, 1704038400, 1704038460},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprint(testCase.input), func(t *testing.T) {
			currentTime := time.Unix(testCase.input, 0)
			currentMinute, nextMinute := utils.GetUnixMinuteRange(currentTime)
			if currentMinute.Unix() != testCase.minStart || nextMinute.Unix() != testCase.minEnd {
				t.Errorf("Error. Expected %v, %v. Got %v, %v", testCase.minStart, testCase.minEnd, currentMinute.Unix(), nextMinute.Unix())
			}
		})
	}
}
