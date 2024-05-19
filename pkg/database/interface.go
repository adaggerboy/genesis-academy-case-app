package database

type IDatabase interface {
	CreateSubscription(email string) (created bool, err error)
	DeleteSubscription(email string) (err error)
	GetSubscriptions() (emails []string, err error)
}
