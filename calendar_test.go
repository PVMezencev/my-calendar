package main

import (
	"testing"
	"time"
)

func TestMonth20240912(t *testing.T) {

	// Инициализируем календарь, в котором первый день недели - понедельник.
	clr := InitCalendar(true)
	// Установим конкретную дату 12 сентября 2024 года (время не принципиально).
	testDay := clr.NewDay(time.Date(2024, 9, 12, 10, 10, 10, 10, time.UTC))
	// Получаем месяц.
	m := testDay.Month()
	// Получаем упорядоченный список недель в виде 2-мерного массива дней.
	weeks := m.MonthOrdered()

	// Первый день заданного месяца должен быть воскресенье, в нашем случае - конец недели.
	firstDay := weeks[0][6]
	if firstDay.Title != "ВС" {
		t.Errorf("название дня %s, а ожидалось ВС", firstDay.Title)
	}

	// Последний день нашего месяца должен быть в понедельник.
	lastDay := m.LastDay()
	if lastDay.Title != "ПН" {
		t.Errorf("название дня %s, а ожидалось ПН", firstDay.Title)
	}

	// Общее количество захватываемых недель - 6, количество полных недель - 4.
	if tc, fc := m.WeeksCount(); tc != 6 {
		t.Errorf("неверное количество недель: %d", tc)
	} else if fc != 4 {
		t.Errorf("неверное количество полных недель: %d", fc)
	}
}
