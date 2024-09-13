package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern

	Применимость:
		- Когда нужно экономить системные ресурсы, повторно используя уже созданные объекты, вместо порождения новых.
		- Возможность пользователям расширять части вашего или библиотеки.
		- Когда нужно делегировать создание объектов: Программа может не знать, какой именно объект нужно создать, а логику создания можно делегировать конкретным подклассам.

	Преимущества:
		- Выделяет код производства продуктов в одно место, упрощая поддержку кода.
		- Добавление новых типов не требует изменения существующего кода, только добавления новых классов и фабрик.
		- Клиентский код не зависит от конкретных типов, он работает только с фабричным методом.
		- В отличии от конструктора позволяет не создавать новый объект.

	Недостатки:
		- Если есть только один тип, фабричный метод может быть излишним.
		- Если типов объектов много, это приведёт к созданию множества фабрик, что усложнит код.

*/

// Notification - интерфейс для всех типов уведомлений
type Notification interface {
	SendNotification(message string) string
}

// EmailNotification - конкретный тип уведомления (Email)
type EmailNotification struct{}

// SendNotification Реализация для Email
func (e EmailNotification) SendNotification(message string) string {
	return "Sending email notification: " + message
}

// SMSNotification - другой конкретный тип уведомления (SMS)
type SMSNotification struct{}

// SendNotification Реализация для SMS
func (s SMSNotification) SendNotification(message string) string {
	return "Sending SMS notification: " + message
}

// NotificationFactory - интерфейс, определяющий фабричный метод
type NotificationFactory interface {
	CreateNotification() Notification
	Send(message string) string
}

// EmailFactory - конкретная фабрика для создания email-уведомлений
type EmailFactory struct{}

func (e EmailFactory) CreateNotification() Notification {
	return EmailNotification{}
}

func (e EmailFactory) Send(message string) string {
	notification := e.CreateNotification()
	return notification.SendNotification(message)
}

// SMSFactory - конкретная фабрика для создания SMS-уведомлений
type SMSFactory struct{}

func (s SMSFactory) CreateNotification() Notification {
	return SMSNotification{}
}

func (s SMSFactory) Send(message string) string {
	notification := s.CreateNotification()
	return notification.SendNotification(message)
}

func main() {
	// Создаем фабрику для email-уведомлений
	emailFactory := EmailFactory{}
	fmt.Println(emailFactory.Send("Hello via Email!")) // "Sending email notification: Hello via Email!"

	// Создаем фабрику для SMS-уведомлений
	smsFactory := SMSFactory{}
	fmt.Println(smsFactory.Send("Hello via SMS!")) // "Sending SMS notification: Hello via SMS!"
}
