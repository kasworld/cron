// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// humantime 기준 정기적으로 반복되는 job들을 관리
package cron

import (
	"fmt"
	"time"
)

const MinuteEnd = 60
const HourEnd = 24
const WeekdayEnd = 7
const MonthdayEnd = 31

type RepeatSkipType uint16

const (
	RST_None RepeatSkipType = iota
	RST_Once
	RST_Many
)

func IsIn(x int, r1, r2 int) bool {
	return x >= r1 && x < r2
}

func CheckHourMinute(hour int, minute int) error {
	if !IsIn(hour, 0, 24) {
		return fmt.Errorf("Invalid Hour %v", hour)
	}
	if !IsIn(minute, 0, 60) {
		return fmt.Errorf("Invalid Minute %v", minute)
	}
	return nil
}

func CalcCommingTime(job IRepeatJob) (time.Time, error) {
	baseTime := time.Now().UTC()
	return job.CalcCommingTimeByBaseTime(baseTime)
}

func CalcPassedTime(job IRepeatJob) (time.Time, error) {
	baseTime := time.Now().UTC()
	return job.CalcPassedTimeByBaseTime(baseTime)
}

func CalcProgressRateByBaseTime(job IRepeatJob, baseTime time.Time) (float64, error) {
	baseTime = baseTime.UTC()

	commingTick, err := job.CalcCommingTimeByBaseTime(baseTime)
	if err != nil {
		return 0, err
	}
	passedTick, err := job.CalcPassedTimeByBaseTime(baseTime)
	if err != nil {
		return 0, err
	}
	passedDur := baseTime.Sub(passedTick)
	passedRate := float64(passedDur) / float64(commingTick.Sub(passedTick))
	return passedRate, nil
}

func CalcProgressRate(job IRepeatJob) (float64, error) {
	baseTime := time.Now().UTC()
	return CalcProgressRateByBaseTime(job, baseTime)
}

func CalcSkipCount(job IRepeatJob, lastTime, currentTime time.Time) (RepeatSkipType, error) {
	lastTime = lastTime.UTC()
	currentTime = currentTime.UTC()

	if lastTime.After(currentTime) {
		return RST_None, fmt.Errorf("lastTime %v > currentTime %v", lastTime, currentTime)
	}
	lastComming, err := job.CalcCommingTimeByBaseTime(lastTime)
	if err != nil {
		return RST_None, err
	}
	if lastComming.After(currentTime) {
		// no cross
		return RST_None, nil
	}
	currentPassed, err := job.CalcPassedTimeByBaseTime(currentTime)
	if err != nil {
		return RST_None, err
	}
	if lastComming == currentPassed {
		return RST_Once, nil
	}
	return RST_Many, nil // or more
}
