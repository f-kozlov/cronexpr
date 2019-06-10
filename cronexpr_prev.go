package cronexpr

/******************************************************************************/

import (
	"sort"
	"time"
)


/******************************************************************************/

func lastOf(slice []int) int {
	return slice[len(slice)-1];
}

func lastDayOfMonth(year int, month time.Month, location *time.Location) int {
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth.Day();
}

// find the index of the value, or the one immediately under; -1 if all elements are greater than value
func findTargetOrLower(slice []int, value int) int {
	i := sort.SearchInts(slice, value)
	if i == 0 && slice[i] != value {
		return -1;
	}
	if i == len(slice) || slice[i] != value {
		i--;
	}
	return i;
}

func (expr *Expression) prevYear(t time.Time) time.Time {

	i := findTargetOrLower(expr.yearList, t.Year()-1);
	if i == -1 {
		return time.Time{}
	}

	// Year changed, need to recalculate actual days of month
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(expr.yearList[i], lastOf(expr.monthList))
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.prevMonth(time.Date(
			expr.yearList[i],
			time.Month(lastOf(expr.monthList)),
			lastDayOfMonth(expr.yearList[i], time.Month(lastOf(expr.monthList)), t.Location()),
			lastOf(expr.hourList),
			lastOf(expr.minuteList),
			lastOf(expr.secondList),
			0,
			t.Location()))
	}
	return time.Date(
		expr.yearList[i],
		time.Month(lastOf(expr.monthList)),
		lastOf(expr.actualDaysOfMonthList),
		lastOf(expr.hourList),
		lastOf(expr.minuteList),
		lastOf(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevMonth(t time.Time) time.Time {

	i := findTargetOrLower(expr.monthList, int(t.Month())-1);

	if i == -1 {
		return expr.prevYear(t)
	}

	// Month changed, need to recalculate actual days of month
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(t.Year(), expr.monthList[i])
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.prevMonth(time.Date(
			t.Year(),
			time.Month(expr.monthList[i]),
			lastDayOfMonth(t.Year(), time.Month(expr.monthList[i]), t.Location()),
			lastOf(expr.hourList),
			lastOf(expr.minuteList),
			lastOf(expr.secondList),
			0,
			t.Location()))
	}

	return time.Date(
		t.Year(),
		time.Month(expr.monthList[i]),
		lastOf(expr.actualDaysOfMonthList),
		lastOf(expr.hourList),
		lastOf(expr.minuteList),
		lastOf(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevDayOfMonth(t time.Time) time.Time {
	// find previous day; if we are out of this month, go to previous month
	i := findTargetOrLower(expr.actualDaysOfMonthList, t.Day()-1);

	if i == -1 {
		return expr.prevMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		expr.actualDaysOfMonthList[i],
		lastOf(expr.hourList),
		lastOf(expr.minuteList),
		lastOf(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevHour(t time.Time) time.Time {
	// Find previous hour; if we are out of this day, go to previous day
	i := findTargetOrLower(expr.hourList, t.Hour()-1);

	if i == -1 {
		return expr.prevDayOfMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		expr.hourList[i],
		lastOf(expr.minuteList),
		lastOf(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevMinute(t time.Time) time.Time {
	// Find previous minute; if we are out of this hour, go to previous one
	i := findTargetOrLower(expr.minuteList, t.Minute()-1);

	if i == -1 {
		return expr.prevHour(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		expr.minuteList[i],
		lastOf(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevSecond(t time.Time) time.Time {

	i := findTargetOrLower(expr.secondList, t.Second()-1);

	if i == -1 {
		return expr.prevMinute(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		expr.secondList[i],
		0,
		t.Location())
}



