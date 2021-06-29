package notify

import "fmt"

type ConsoleSender struct{}

func (cs *ConsoleSender) Send(text string) {
	fmt.Println(text)
}
