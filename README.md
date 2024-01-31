## Календарь на Golang.

#### Установка:
```bash
go get github.com/PVMezencev/my-calendar
```

#### Использование:
```go
package main

import (
	calendar "github.com/PVMezencev/my-calendar"
	"log"
)

func main() {
	// Инициализируем календарь, в котором первый день недели - понедельник.
	clr := calendar.InitCalendar(true)
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
```