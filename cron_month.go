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

var _ IRepeatJob = &MonthlyRepeatJob{}

// humantime 윌 반복
type MonthlyRepeatJob struct {
	Name   string
	Day    []bool // 1(0) ~ 31(30)
	Hour   int    // 0 ~ 23
	Minute int    // 0 ~ 59
}

func NewMonthlyRepeatJob(name string, hour int, minute int, scheduleList []bool) (*MonthlyRepeatJob, error) {
	if err := CheckHourMinute(hour, minute); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if l := len(scheduleList); l != MonthdayEnd {
		return nil, fmt.Errorf("Fail To NewMonthlyRepeatJob, Incorrect List Len, %v != %v", l, MonthdayEnd)
	}

	return &MonthlyRepeatJob{
		Name:   name,
		Day:    scheduleList,
		Hour:   hour,
		Minute: minute,
	}, nil
}

func (mrj MonthlyRepeatJob) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "%vMonthlyRepeatJob[BaseTime %v:%v:00 Day:",
		mrj.Name, mrj.Hour, mrj.Minute)

	for i, v := range mrj.Day {
		if v == true {
			fmt.Fprintf(&buff, "%v,", i)
		}
	}

	fmt.Fprintf(&buff, "]")
	return buff.String()
}

func (mrj *MonthlyRepeatJob) SetSchedule(scheduleList []bool) error {
	if len(scheduleList) != len(mrj.Day) {
		return fmt.Errorf("Incorrect List Len, %v != %v", len(scheduleList), len(mrj.Day))
	}

	mrj.Day = scheduleList

	return nil
}

func (mrj MonthlyRepeatJob) CalcCommingTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= MonthdayEnd; i++ {
		dstTime := baseTime.AddDate(0, 0, i)
		if mrj.Day[dstTime.Day()-1] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				mrj.Hour, mrj.Minute, 0, 0, time.UTC)
			if t.After(baseTime) {
				// now find comming tick
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("comming tick not found %v", mrj)
}

func (mrj MonthlyRepeatJob) CalcPassedTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= MonthdayEnd; i++ {
		dstTime := baseTime.AddDate(0, 0, -i)
		// log.Printf("%v %v", dstTime, dstTime.Day())
		if mrj.Day[dstTime.Day()-1] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				mrj.Hour, mrj.Minute, 0, 0, time.UTC)
			if t.Before(baseTime) {
				// now find passed tick
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("passed tick not found %+v", mrj)
}
