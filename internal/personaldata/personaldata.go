// Package personaldata хранит и выводит персональные данные пользователя.
package personaldata

import "fmt"

// Personal содержит имя, вес в килограммах и рост в метрах.
type Personal struct {
	Name   string
	Weight float64
	Height float64
}

// Print выводит персональные данные в стандартный поток вывода.
func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f\nРост: %.2f\n", p.Name, p.Weight, p.Height)
}
