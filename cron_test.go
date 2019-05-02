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
	"testing"
	"time"
)

func TestNewHourly(t *testing.T) {
	now := time.Now()

	hourlySchedule := make([]bool, HourEnd)
	for i, _ := range hourlySchedule {
		hourlySchedule[i] = true
	}
	rj, err := NewHourlyRepeatJob("hourly", hourlySchedule, 25)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pt, err := rj.CalcPassedTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	ct, err := rj.CalcCommingTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pr, err := CalcProgressRate(rj)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v\n%v\n%v\n%v\n%v", rj, pt, now, ct, pr)
}

func TestNewDaily(t *testing.T) {
	now := time.Now()

	dailySchedule := make([]bool, WeekdayEnd)
	for i, _ := range dailySchedule {
		dailySchedule[i] = true
	}

	rj, err := NewWeeklyRepeatJob("Daily", 0, 0, dailySchedule)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pt, err := rj.CalcPassedTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	ct, err := rj.CalcCommingTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pr, err := CalcProgressRate(rj)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v\n%v\n%v\n%v\n%v", rj, pt, now, ct, pr)

	from := now.Add(-24 * 3600 * time.Second)
	to := now.Add(24 * 3600 * time.Second)
	sc, err := CalcSkipCount(rj, from, to)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v\n%v\n%v\n%v",
		now,
		from,
		to,
		sc,
	)

}

func TestNewWeekly(t *testing.T) {
	now := time.Now()

	weeklySchedule := make([]bool, WeekdayEnd)
	weeklySchedule[0] = true

	rj, err := NewWeeklyRepeatJob("Weekly", 0, 0, weeklySchedule)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pt, err := rj.CalcPassedTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	ct, err := rj.CalcCommingTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pr, err := CalcProgressRate(rj)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v\n%v\n%v\n%v\n%v", rj, pt, now, ct, pr)
}

func TestNewMonthly(t *testing.T) {
	now := time.Now()

	monthlySchedule := make([]bool, MonthdayEnd)
	monthlySchedule[0] = true
	rj, err := NewMonthlyRepeatJob("Monthly", 0, 0, monthlySchedule)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pt, err := rj.CalcPassedTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	ct, err := rj.CalcCommingTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pr, err := CalcProgressRate(rj)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v\n%v\n%v\n%v\n%v", rj, pt, now, ct, pr)

	dCron, err := NewMonthlyRepeatJob("Monthly", 0, 0, monthlySchedule)
	if err != nil {
		t.Fatalf("%v", err)
	}
	dCron.Day[0] = true
	t.Logf("%v", dCron)
	dtick, err := dCron.CalcCommingTimeByBaseTime(now)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", dtick)
}
