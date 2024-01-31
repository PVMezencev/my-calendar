package main

import (
	"fmt"
	"log"
	"time"

	calendar "github.com/PVMezencev/my-calendar"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramCalendar Использование календаря для генерации callback-кнопок в формате месяца.
func TelegramCalendar(bot *tgbotapi.BotAPI,
	msg *tgbotapi.Message,
	date time.Time,
	caption, cbCommand string,
	editMessage bool) {
	if bot == nil || msg == nil {
		return
	}
	// Инициализировать календарь.
	clndr := calendar.InitCalendar(false)
	// Инициализировать тип Месяц календаря.
	indexDay := clndr.NewDay(date)
	month := indexDay.Month()
	weekTitle := month.WeekTitle()

	// Заполнить строку с названием календаря.
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	monthTitleBtn := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("cl %s %s",
			cbCommand,
			month.Previous().LastDay().Time.Format("2006-01-02"))),
		tgbotapi.NewInlineKeyboardButtonData(month.MonthTitleWithYear(), month.MonthTitleWithYear()),
		tgbotapi.NewInlineKeyboardButtonData(">", fmt.Sprintf("cl %s %s",
			cbCommand,
			month.Next().FirstDay().Time.Format("2006-01-02"))),
	}
	rows = append(rows, monthTitleBtn)

	// Заполнить строку дней недели.
	weekTitleBtn := []tgbotapi.InlineKeyboardButton{}
	for _, wt := range weekTitle {
		wtBtn := tgbotapi.NewInlineKeyboardButtonData(wt, wt)
		weekTitleBtn = append(weekTitleBtn, wtBtn)
	}
	rows = append(rows, weekTitleBtn)

	// Заполнить числа календаря.
	orderedDays := month.MonthOrdered()
	for _, week := range orderedDays {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		for _, day := range week {
			btnText := " "
			clbData := "-"
			if day != nil {
				btnText = fmt.Sprintf("%d", day.Time.Day())
				clbData = fmt.Sprintf("%s %s", cbCommand, day.Time.Format("2006-01-02"))
				if day.Time.Format("2006-01-02") == time.Now().UTC().Format("2006-01-02") {
					btnText = fmt.Sprintf("%d%s", day.Time.Day(), "\U0001F7E2")
				}
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(btnText, clbData)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}
	// Сформируем сообщение.
	docConf := tgbotapi.NewMessage(msg.Chat.ID, caption)
	if editMessage {
		editConf := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, caption)
		kbmk := tgbotapi.NewInlineKeyboardMarkup(rows...)
		editConf.ReplyMarkup = &kbmk
		_, err := bot.Send(editConf)
		if err != nil {
			log.Printf("TelegramCalendar.bot.Send(editMessage %s): %s", caption, err)
		}
		return
	}
	docConf.ParseMode = tgbotapi.ModeMarkdown
	docConf.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	_, err := bot.Send(docConf)
	if err != nil {
		log.Printf("TelegramCalendar.bot.Send(%s): %s", caption, err)
	}
}
