// Package actioninfo обрабатывает наборы записей через общий интерфейс DataParser.
package actioninfo

import (
	"fmt"
	"log"
)

// DataParser объединяет типы, которые умеют разбирать запись и возвращать отчёт.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info последовательно разбирает записи и выводит успешные отчёты.
// Ошибка одной записи логируется и не прерывает обработку остальных записей.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Printf("ошибка разбора строки %q: %v", data, err)
			continue
		}

		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("ошибка формирования отчёта для строки %q: %v", data, err)
			continue
		}

		// Методы ActionInfo возвращают строку с завершающим переводом строки.
		fmt.Print(info)
	}
}
