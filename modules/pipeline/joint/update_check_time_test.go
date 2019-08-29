/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package joint

import (
	"fmt"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xirtah/gopa-framework/core/model"
)

func TestInitGrabVelocityArr(t *testing.T) {
	steps := initFetchRateArr("24h,12h,6h,3h,1h30m,45m,20m,10m,1m")
	fmt.Println(steps)
	assert.Equal(t, steps, []int{86400, 43200, 21600, 10800, 5400, 2700, 1200, 600, 60})
}

func TestSetSnapNextCheckTime(t *testing.T) {
	steps := initFetchRateArr("10m,5m,3m,2m,1m")
	startStep := initFetchRateArr("3m")[0]

	fmt.Println("steps,", steps)

	zeroTime := time.Time{} //instantiate a zero time i.e 0001-01-01 00:00:00 +0000 UTC
	startTime :=  zeroTime.Add(1 * oneSecond) 
	oneSecond, _ := time.ParseDuration("1s")
	oneMinute, _ := time.ParseDuration("1m")
	
	//no last check, no next check
	fmt.Println("no last check, no next check - with change")
	currentTime := startTime.Add(1 * oneSecond)
	context := model.Context{}
	context.Set(model.CONTEXT_TASK_LastCheck, zeroTime)
	context.Set(model.CONTEXT_TASK_NextCheck, zeroTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 := context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 := context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval := getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 180, timeInterval)

	fmt.Println()
	
	fmt.Println("no last check, no next check - with no change")
	currentTime = startTime.Add(1 * oneSecond)
	context = model.Context{}
	context.Set(model.CONTEXT_TASK_LastCheck, zeroTime)
	context.Set(model.CONTEXT_TASK_NextCheck, zeroTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 180, timeInterval)

	fmt.Println()
	
	//no last check, yes next check
	fmt.Println("no last check, no next check - with change")
	currentTime = startTime.Add(1 * oneSecond)
	context = model.Context{}
	context.Set(model.CONTEXT_TASK_LastCheck, zeroTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 180, timeInterval)

	fmt.Println()
	
	fmt.Println("no last check, no next check - with no change")
	currentTime = startTime.Add(1 * oneSecond)
	context = model.Context{}
	context.Set(model.CONTEXT_TASK_LastCheck, zeroTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 180, timeInterval)

	fmt.Println()

	//following test cases are for the different variations of yes last check, yes next check
	fmt.Println("update 1s with no change")
	context = model.Context{}
	currentTime = startTime.Add(1 * oneSecond)
	fmt.Println(startTime)
	fmt.Println(!(startTime.IsZero()))
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println()

	fmt.Println("update 10m with no change")
	currentTime = startTime.Add(10 * oneMinute)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println()

	fmt.Println("update 20m with no change")
	context = model.Context{}
	currentTime = startTime.Add(10 * oneMinute)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println()

	fmt.Println("update 2m with change")
	currentTime = startTime.Add(120 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println("update 10s with change")
	currentTime = startTime.Add(10 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println("update 1000s with change")
	currentTime = startTime.Add(1000 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println("update 500s with change")
	currentTime = startTime.Add(500 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 300, timeInterval)

	fmt.Println("update 600s with change")
	currentTime = startTime.Add(600 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, startTime)
	context.Set(model.CONTEXT_TASK_NextCheck, currentTime)
	updateNextCheckTime(&context, currentTime, startStep, steps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 300, timeInterval)
}
