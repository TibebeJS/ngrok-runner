package notifier

type Notifier interface {
	Notify(message string) (bool, error)
}
