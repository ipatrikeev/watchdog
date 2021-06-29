package notify

type Sender interface {
	Send(text string)
}
