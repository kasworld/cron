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

var _ IRepeatJob = &HourlyRepeatJob{}

// humantime 시각 기준 반복
type HourlyRepeatJob struct {
	Name   string
	Hour   []bool // 0 ~ 23
	Minute int    // 0 ~ 59
}

func NewHourlyRepeatJob(name string, scheduleList []bool, minute int) (*HourlyRepeatJob, error) {
	if l := len(scheduleList); l != HourEnd {
		return nil, fmt.Errorf("Fail To NewHourlyRepeatJob, Incorrect List Len, %v != %v", l, MonthdayEnd)
	}
	return &HourlyRepeatJob{
		Name:   name,
		Hour:   scheduleList,
		Minute: minute,
	}, nil
}

func (hrj HourlyRepeatJob) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff,
		"%vHourlyRepeatJob[Minute:%v Hour:",
		hrj.Name, hrj.Minute)

	for i, v := range hrj.Hour {
		if v == true {
			fmt.Fprintf(&buff, "%v,", i)
		}
	}

	fmt.Fprintf(&buff, "]")
	return buff.String()
}

func (hrj *HourlyRepeatJob) SetSchedule(scheduleList []bool) error {
	if len(scheduleList) != len(hrj.Hour) {
		return fmt.Errorf("Incorrect List Len, %v != %v", len(scheduleList), len(hrj.Hour))
	}

	hrj.Hour = scheduleList

	return nil
}

func (hrj HourlyRepeatJob) CalcCommingTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= HourEnd; i++ {
		dur := time.Duration(time.Hour) * time.Duration(i)
		dstTime := baseTime.Add(dur)

		if hrj.Hour[dstTime.Hour()] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				dstTime.Hour(), hrj.Minute, 0, 0, time.UTC)
			if t.After(baseTime) {
				// now find comming tick
				return t, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("comming tick not found %v", hrj)
}

func (hrj HourlyRepeatJob) CalcPassedTimeByBaseTime(baseTime time.Time) (time.Time, error) {
	baseTime = baseTime.UTC()
	for i := 0; i <= HourEnd; i++ {
		dur := time.Duration(time.Hour) * time.Duration(i)
		dstTime := baseTime.Add(-dur)
		if hrj.Hour[dstTime.Hour()] {
			t := time.Date(dstTime.Year(), dstTime.Month(), dstTime.Day(),
				dstTime.Hour(), hrj.Minute, 0, 0, time.UTC)

			if t.Before(baseTime) {
				// now find passed tick
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("passed tick not found %+v", hrj)
}
