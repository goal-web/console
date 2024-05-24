package scheduling

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/golang-module/carbon/v2"
	"strconv"
	"strings"
	"time"
)

func NewEvent(mutex *Mutex, callback any, timezone string) *Event {
	return &Event{
		callback:           callback,
		mutex:              mutex,
		filters:            make([]filter, 0),
		rejects:            make([]filter, 0),
		beforeCallbacks:    make([]func(), 0),
		afterCallbacks:     make([]func(), 0),
		withoutOverlapping: false,
		onOneServer:        false,
		timezone:           timezone,
		expression:         "0 * * * * * *",
		mutexName:          "",
		expiresAt:          0,
	}
}

type filter func() bool

type Event struct {
	callback any

	mutex           *Mutex
	filters         []filter
	rejects         []filter
	beforeCallbacks []func()
	afterCallbacks  []func()

	withoutOverlapping bool
	onOneServer        bool

	timezone   string
	expression string
	mutexName  string
	expiresAt  time.Duration
}

func (event *Event) Years(years ...string) contracts.ScheduleEvent {
	if len(years) > 0 {
		return event.SpliceIntoPosition(6, strings.Join(years, ","))
	}
	return event
}

func (event *Event) Expression() string {
	return event.expression
}

func (event *Event) EveryThirtySeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "0,30")
}

func (event *Event) EveryFifteenSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/15")
}

func (event *Event) EveryTenSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/10")
}

func (event *Event) EveryFiveSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/5")
}

func (event *Event) EveryFourSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/4")
}

func (event *Event) EveryThreeSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/3")
}

func (event *Event) EveryTwoSeconds() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*/2")
}

func (event *Event) EverySecond() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(0, "*")
}

func (event *Event) WithoutOverlapping(expiresAt int) contracts.ScheduleEvent {
	event.expiresAt = time.Duration(expiresAt) * time.Second
	event.withoutOverlapping = true
	return event.Skip(func() bool {
		return event.mutex.Exists(event)
	})
}

func (event *Event) Run(application contracts.Application) {
	if !event.FiltersPass() {
		return
	}
	defer event.removeMutex()
	if event.withoutOverlapping && !event.mutex.Create(event) {
		return
	}
	application.Call(event.callback)
}

func (event *Event) removeMutex() {
	if event.withoutOverlapping {
		event.mutex.Forget(event)
	}
}
func (event *Event) OnOneServer() contracts.ScheduleEvent {
	event.onOneServer = true
	return event
}

func (event *Event) Timezone(timezone string) contracts.ScheduleEvent {
	event.timezone = timezone
	return event
}

func (event *Event) Days(day string, days ...string) contracts.ScheduleEvent {
	days = append([]string{day}, days...)
	return event.SpliceIntoPosition(5, strings.Join(days, ","))
}

func (event *Event) YearlyOn(month time.Month, dayOfMonth int, timeStr string) contracts.ScheduleEvent {
	event.DailyAt(timeStr)

	return event.SpliceIntoPosition(3, strconv.Itoa(dayOfMonth)).
		SpliceIntoPosition(4, strconv.Itoa(int(month)))
}

func (event *Event) Yearly() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "0").
		SpliceIntoPosition(3, "1").
		SpliceIntoPosition(4, "1")
}

func (event *Event) Quarterly() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "0").
		SpliceIntoPosition(3, "1").
		SpliceIntoPosition(4, "1-12/3")
}

func (event *Event) LastDayOfMonth(timeStr string) contracts.ScheduleEvent {
	event.DailyAt(timeStr)

	return event.When(func() bool {
		return carbon.Now(event.timezone).Day() == carbon.Now(event.timezone).EndOfMonth().Day()
	})
}

func (event *Event) TwiceMonthly(first, second int, timeStr string) contracts.ScheduleEvent {
	event.DailyAt(timeStr)
	return event.SpliceIntoPosition(3, fmt.Sprintf("%d,%d", first, second))
}

func (event *Event) MonthlyOn(dayOfMonth int, timeStr string) contracts.ScheduleEvent {
	event.DailyAt(timeStr)
	return event.SpliceIntoPosition(3, strconv.Itoa(dayOfMonth))
}

func (event *Event) Monthly() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "0").
		SpliceIntoPosition(3, "1")
}

func (event *Event) WeeklyOn(dayOfWeek time.Weekday, timeStr string) contracts.ScheduleEvent {
	return event.DailyAt(timeStr).Days(strconv.Itoa(int(dayOfWeek)))
}

func (event *Event) Weekly() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "0").
		SpliceIntoPosition(5, "0")
}

func (event *Event) Sundays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Sunday))
}

func (event *Event) Saturdays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Saturday))
}

func (event *Event) Fridays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Friday))
}

func (event *Event) Thursdays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Thursday))
}

func (event *Event) Wednesdays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Wednesday))
}

func (event *Event) Tuesdays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Tuesday))
}

func (event *Event) Mondays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d", time.Monday))
}

func (event *Event) Weekends() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d,%d", time.Saturday, time.Sunday))
}

func (event *Event) Weekdays() contracts.ScheduleEvent {
	return event.Days(fmt.Sprintf("%d-%d", time.Monday, time.Friday))
}

func (event *Event) TwiceDailyAt(first, second, offset int) contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, strconv.Itoa(offset)).
		SpliceIntoPosition(2, fmt.Sprintf("%d,%d", first, second))
}

func (event *Event) TwiceDaily(first, second int) contracts.ScheduleEvent {
	return event.TwiceDailyAt(first, second, 0)
}

func (event *Event) DailyAt(timeStr string) contracts.ScheduleEvent {
	segments := strings.Split(timeStr, ":")
	event.SpliceIntoPosition(2, segments[0])

	if len(segments) == 2 {
		return event.SpliceIntoPosition(1, segments[1])
	} else {
		return event.SpliceIntoPosition(1, "0")
	}
}

func (event *Event) Daily() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "0")
}

func (event *Event) EverySixHours() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "*/6")
}

func (event *Event) EveryFourHours() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "*/4")
}

func (event *Event) EveryThreeHours() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "*/3")
}

func (event *Event) EveryTwoHours() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0").
		SpliceIntoPosition(2, "*/2")
}

func (event *Event) HourlyAt(offset ...int) contracts.ScheduleEvent {
	offsetStrings := make([]string, 0)
	for _, offsetInt := range offset {
		offsetStrings = append(offsetStrings, strconv.Itoa(offsetInt))
	}
	return event.SpliceIntoPosition(1, strings.Join(offsetStrings, ","))
}

func (event *Event) Hourly() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0")
}

func (event *Event) EveryThirtyMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "0,30")
}

func (event *Event) EveryFifteenMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/15")
}

func (event *Event) EveryTenMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/10")
}

func (event *Event) EveryFiveMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/5")
}

func (event *Event) EveryFourMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/4")
}

func (event *Event) EveryThreeMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/3")
}

func (event *Event) EveryTwoMinutes() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*/2")
}

func (event *Event) EveryMinute() contracts.ScheduleEvent {
	return event.SpliceIntoPosition(1, "*")
}

func (event *Event) FiltersPass() bool {
	for _, filter := range event.filters {
		if !filter() {
			return false
		}
	}
	for _, reject := range event.rejects {
		if reject() {
			return false
		}
	}
	return true
}
func (event *Event) When(filter func() bool) contracts.ScheduleEvent {
	event.filters = append(event.filters, filter)
	return event
}
func (event *Event) Skip(reject func() bool) contracts.ScheduleEvent {
	event.rejects = append(event.rejects, reject)
	return event
}

func (event *Event) Cron(expression string) contracts.ScheduleEvent {
	event.expression = expression
	return event
}

func (event *Event) Between(startTime, endTimeStr string) contracts.ScheduleEvent {
	return event.When(event.inTimeInterval(startTime, endTimeStr))
}

func (event *Event) UnlessBetween(startTime, endTimeStr string) contracts.ScheduleEvent {
	return event.Skip(event.inTimeInterval(startTime, endTimeStr))
}

func (event *Event) inTimeInterval(startTime, endTimeStr string) func() bool {
	var (
		startAt = carbon.Now().ParseByFormat(startTime, "H:i", event.timezone)
		endAt   = carbon.Now().ParseByFormat(endTimeStr, "H:i", event.timezone)
	)

	if endAt.Lt(startAt) {
		if startAt.Gt(carbon.Now(event.timezone).SetYear(0000).SetMonth(1).SetDay(1)) {
			startAt.SubDay()
		} else {
			endAt.AddDay()
		}
	}

	return func() bool {
		now := carbon.Now(event.timezone).SetYear(0000).SetMonth(1).SetDay(1)
		return now.Between(startAt, endAt)
	}
}

func (event *Event) MutexName() string {
	return event.mutexName
}

func (event *Event) SetMutexName(mutexName string) contracts.ScheduleEvent {
	event.mutexName = mutexName
	return event
}

func (event *Event) SpliceIntoPosition(position int, value string) contracts.ScheduleEvent {
	segments := strings.Split(event.expression, " ")
	segments[position] = value
	return event.Cron(strings.Join(segments, " "))
}
