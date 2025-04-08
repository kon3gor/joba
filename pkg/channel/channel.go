package channel

type C interface {
	SendMessage(string) error
}
