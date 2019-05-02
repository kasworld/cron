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

package cron

import (
	"bytes"
	"fmt"
	"time"
)

const WeekOn = "SMTWTFS"
const WeekOff = "smtwtfs"

var _ IRepeatJob = &WeeklyRepeatJob{}

// humantime 주 반복
type WeeklyRepeatJob struct {
	Name      string
	DayOfWeek []bool // Sun(0) ~ Sat(6)
	Hour      int    // 0 ~ 23
	Minute    int    // 0 ~ 59
}

func NewWeeklyRepeatJob(name string, hour int, minute int, scheduleList []bool) (*WeeklyRepeatJob, error) {
	if err := CheckHourMinute(hour, minute); err != nil {

		return nil, fmt.Errorf("%v", err)
	}

	if l := len(scheduleList); l != WeekdayEnd {

		return nil, fmt.Errorf("Fail to NewWeeklyRepeatJob, Incorrect List Len, %v != %v", l, WeekdayEnd)
	}

	return &WeeklyRepeatJob{
		Name:      name,
		DayOfWeek: scheduleList,
		Hour:      hour,
		Minute:    minute,
	}, nil
}

func (wrj WeeklyRepeatJob) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "%vWeeklyRepeatJob[BaseTime %v:%v:00 Day:",
		wrj.Name, wrj.Hour, wrj.Minute)

	for i, v := range wrj.DayOfWeek {
		if v {
			fmt.Fprintf(&buff, "%c", WeekOn[i])
		} else {
			fmt.Fprintf(&buff, "%c", WeekOff[i])
		}
	}
	fmt.Fprintf(&buff, "]")
	return buff.String()
}

func (wrj *WeeklyRepeatJob) SetSchedule(scheduleList []bool) error {
	if len(scheduleList) != len(wrj.DayOfWeek) {
		return fmt.Errorf("Incorrect List Len, %v != %v", len(scheduleList), len(wrj.DayOfWeek))
	}

	wrj.DayOfWeek = scheduleList

	return nil
}

func (wrj WeeklyRepeatJob) CalcCommingTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= WeekdayEnd; i++ {
		dstTime := baseTime.AddDate(0, 0, i)
		if wrj.DayOfWeek[dstTime.Weekday()] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				wrj.Hour, wrj.Minute, 0, 0, time.UTC)
			if t.After(baseTime) {
				// now find comming tick
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("comming tick not found %v", wrj)
}

func (wrj WeeklyRepeatJob) CalcPassedTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= WeekdayEnd; i++ {
		dstTime := baseTime.AddDate(0, 0, -i)
		if wrj.DayOfWeek[dstTime.Weekday()] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				wrj.Hour, wrj.Minute, 0, 0, time.UTC)
			if t.Before(baseTime) {
				// now find passed tick
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("passed tick not found %+v", wrj)
}
