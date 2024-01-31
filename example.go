package main

import "log"

func main() {
	// Инициализируем календарь, в котором первый день недели - понедельник.
	clr := InitCalendar(true)
	// Получаем сегодня.
	today := clr.Today()
	// Получаем месяц.
	m := today.Month()

	// Распечатаем дни месяца по неделям.
	log.Printf("%s", m.MonthTitleWithYear())
	for _, week := range m.MonthOrdered() {
		for _, day := range week {
			if day != nil {
				log.Printf(`%s: %s`, day.Title, day.Time.Format("02.01.2006"))
			} else {
				log.Printf("-")
			}
		}
		log.Printf("----------")
	}
}
