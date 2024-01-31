package my_calendar

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

var WeekDaysShortTile = map[int]string{
	0: "ВС", 1: "ПН", 2: "ВТ", 3: "СР", 4: "ЧТ", 5: "ПТ", 6: "СБ",
}

var MonthsTile = map[int]string{
	1: "Январь", 2: "Февраль", 3: "Март", 4: "Апрель", 5: "Май", 6: "Июнь",
	7: "Июль", 8: "Август", 9: "Сентябрь", 10: "Октябрь", 11: "Ноябрь", 12: "Декабрь",
}

type Day struct {
	calendar *Calendar
	Time     time.Time `json:"time"`
	Title    string    `json:"title"`
}

func (d Day) StartOfWeek() time.Time {
	offset := int(d.calendar.weekDays[0] - d.Time.Weekday())
	if offset > 0 {
		offset = -6
	}
	start := d.Time.AddDate(0, 0, offset)
	return start
}

func (d Day) Week() Week {
	week := Week{
		calendar: d.calendar,
		Days:     make(map[time.Weekday]*Day, 7),
	}

	previousWeekend := d.StartOfWeek().Add(time.Duration(-24) * time.Hour)
	for idx, dw := range d.calendar.weekDays {
		nd := d.calendar.NewDay(previousWeekend.Add(time.Duration((idx+1)*24) * time.Hour))
		week.Days[dw] = &nd
	}

	return week
}

func (d Day) Month() Month {
	month := Month{
		calendar: d.calendar,
	}

	// Первый день текущего месяца.
	firstDay := time.Date(d.Time.Year(), d.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Последний день текущего месяца.
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)

	checkDay := firstDay
	// Заполним месяц.
	for lastDay.After(checkDay) {
		week := d.calendar.NewDay(checkDay).Week()
		month.Weeks = append(month.Weeks, week)
		// Следующий день возьмем из начала созданной недели и прибавим ему 7 дней.
		startWeek := week.DayByIndex(d.calendar.weekDays[0])
		checkDay = startWeek.Time.Add(time.Duration(7*24) * time.Hour)
	}

	month.Title = month.MonthTitle()

	return month
}

type Week struct {
	calendar *Calendar
	Days     map[time.Weekday]*Day `json:"days"`
}

func (w Week) DayByIndex(index time.Weekday) *Day {
	if d, ok := w.Days[index]; ok {
		return d
	}
	return nil
}

func (w Week) isEmpty() bool {
	if w.calendar.weekStartMonday {
		return w.Days[time.Monday] == nil && w.Days[time.Sunday] == nil
	}
	return w.Days[time.Sunday] == nil && w.Days[time.Saturday] == nil
}

type Month struct {
	calendar *Calendar
	Title    string `json:"title"`
	Weeks    []Week `json:"weeks"`
}

func (m Month) TimeMonth() time.Month {
	return m.timeMonth()
}

func (m Month) timeMonth() time.Month {
	day := m.anyDay()
	return day.Time.Month()
}

func (m Month) anyDay() Day {
	// У месяца точно есть 2 неделя.
	week := m.Weeks[1]
	// Во второй неделе точно есть понедельник.
	return *week.Days[time.Monday]
}

func (m Month) MonthOrdered() [][]*Day {
	monthDays := make([][]*Day, 0)
	for _, week := range m.Weeks {
		weekDays := make([]*Day, 0)
		for _, twk := range m.calendar.weekDays {
			if d, ok := week.Days[twk]; ok && d != nil && d.Time.Month() == m.TimeMonth() {
				weekDays = append(weekDays, d)
			} else {
				weekDays = append(weekDays, nil)
			}
		}
		monthDays = append(monthDays, weekDays)
	}
	return monthDays
}

func (m Month) MonthTitle() string {
	return MonthsTile[int(m.TimeMonth())]
}

func (m Month) MonthTitleWithYear() string {
	return fmt.Sprintf("%s, %d", MonthsTile[int(m.TimeMonth())], m.anyDay().Time.Year())
}

func (m Month) MonthHTML() template.HTML {

	var text string
	text += `<table>`
	text += fmt.Sprintf(`<caption>%s</caption>`, m.MonthTitle())
	text += `<tbody>`
	text += `<tr>`

	weekTitle := m.WeekTitle()
	args := make([]interface{}, 0)
	for _, wt := range weekTitle {
		args = append(args, wt)
	}
	text += fmt.Sprintf(strings.Repeat("<th>%s</th>", len(weekTitle)), args...)

	text += `</tr>`

	ordered := m.MonthOrdered()
	for _, week := range ordered {
		text += `<tr>`
		for _, day := range week {
			if day != nil {
				text += fmt.Sprintf(`<td id= "%s">%d</td>`, day.Time.Format("02.01.2006"), day.Time.Day())
			} else {
				text += `<td>&nbsp;</td>`
			}
		}
		text += `</tr>`
	}

	text += `</tbody>`
	text += `</table>`
	return template.HTML(text)
}

func (m Month) WeekTitle() []string {
	titles := make([]string, 0)
	for _, v := range m.calendar.weekDays {
		titles = append(titles, WeekDaysShortTile[int(v)])
	}
	return titles[:]
}

func (m Month) FirstDay() Day {
	day := m.anyDay()
	// Первый день текущего месяца.
	firstDay := time.Date(day.Time.Year(), day.Time.Month(), 1, 0, 0, 0, 0, time.UTC)

	return m.calendar.NewDay(firstDay)
}

func (m Month) LastDay() Day {
	day := m.anyDay()
	// Первый день текущего месяца.
	firstDay := time.Date(day.Time.Year(), day.Time.Month(), 1, 0, 0, 0, 0, time.UTC)
	// Последний день текущего месяца.
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return m.calendar.NewDay(lastDay)
}

func (m Month) Next() Month {
	// Последний день текущего месяца.
	lastDay := m.LastDay()
	// Получим день следующего месяца.
	date := lastDay.Time.Add(time.Duration(24) * time.Hour)

	newDay := m.calendar.NewDay(date)
	return newDay.Month()
}

func (m Month) Previous() Month {
	// Первый день текущего месяца.
	firstDay := m.FirstDay()
	// Получим день предыдущего месяца.
	date := firstDay.Time.Add(time.Duration(-24) * time.Hour)

	newDay := m.calendar.NewDay(date)
	return newDay.Month()
}

func (m Month) weekIsNotMine(w Week) bool {
	if w.isEmpty() {
		return true
	}
	if m.calendar.weekStartMonday {
		return (w.Days[time.Monday] == nil || w.Days[time.Monday].Time.Month() != m.TimeMonth()) &&
			(w.Days[time.Sunday] == nil || w.Days[time.Sunday].Time.Month() != m.TimeMonth())
	}
	return (w.Days[time.Sunday] == nil || w.Days[time.Sunday].Time.Month() != m.TimeMonth()) &&
		(w.Days[time.Saturday] == nil || w.Days[time.Saturday].Time.Month() != m.TimeMonth())
}

func (m Month) WeeksCount() (int, int) {
	total := 0
	notFull := 0
	for _, week := range m.MonthOrdered() {
		total++
		for _, day := range week {
			if day == nil {
				notFull++
				break
			}
		}
	}
	return total, total - notFull
}

type Calendar struct {
	weekStartMonday bool
	weekDays        []time.Weekday
}

func (c Calendar) NewDay(t time.Time) Day {
	if t.IsZero() {
		t = time.Now()
	}
	day := Day{
		calendar: &c,
		Time:     t,
	}
	day.Title = WeekDaysShortTile[int(day.Time.Weekday())]
	return day
}

func (c Calendar) Today() Day {
	day := c.NewDay(time.Time{})
	day.Title = WeekDaysShortTile[int(day.Time.Weekday())]
	return day
}

func (c Calendar) Yesterday() Day {
	day := c.Today()
	day.Time.Add(time.Duration(-24) * time.Hour)
	day.Title = WeekDaysShortTile[int(day.Time.Weekday())]
	return day
}

func (c Calendar) Tomorrow() Day {
	day := c.Today()
	day.Time.Add(time.Duration(24) * time.Hour)
	day.Title = WeekDaysShortTile[int(day.Time.Weekday())]
	return day
}

func InitCalendar(weekStartMonday bool) Calendar {
	cl := Calendar{
		weekStartMonday: weekStartMonday,
	}
	if weekStartMonday {
		cl.weekDays = []time.Weekday{
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
			time.Sunday,
		}
	} else {
		cl.weekDays = []time.Weekday{
			time.Sunday,
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
		}
	}
	return cl
}
