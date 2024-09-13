package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

	Применимость:
		- Когда важно, чтобы обработчики выполнялись один за другим в строгом порядке.
		- Когда набор объектов, способных обработать запрос, должен задаваться динамически.
		- Когда программа должна обрабатывать разнообразные запросы несколькими способами, но заранее неизвестно, какие конкретно запросы будут приходить и какие обработчики для них понадобятся.

	Преимущества:
		- Отправитель запроса не знает, кто его обработает, что делает систему более модульной.
		- Возможность добавлять, удалять или изменять порядок обработчиков в цепочке без изменения отправителей запросов.
		- Каждый обработчик отвечает за свою часть обработки, что минимизирует дублирование кода.

	Недостатки:
		- Запрос может остаться никем не обработанным.
		- Когда цепочка длинная или плохо организована, запрос пройдёт через множество обработчиков, прежде чем будет обработан.
		- Когда не ясно, какой обработчик должен взять на себя обработку запроса, это может усложнить диагностику проблем.

*/

// Handler Определяем интерфейс для всех обработчиков
type Handler interface {
	SetNext(handler Handler)
	Handle(request string)
}

// BaseHandler Базовая структура обработчика
type BaseHandler struct {
	next Handler
}

// SetNext Реализация установки следующего обработчика
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// Handle Реализация передачи запроса следующему обработчику
func (h *BaseHandler) Handle(request string) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

// HandlerOne Конкретный обработчик 1
type HandlerOne struct {
	BaseHandler
}

func (h *HandlerOne) Handle(request string) {
	if request == "One" {
		fmt.Println("HandlerOne обработал запрос")
	} else {
		fmt.Println("HandlerOne не смог обработать запрос, передача дальше")
		h.BaseHandler.Handle(request)
	}
}

// HandlerTwo Конкретный обработчик 2
type HandlerTwo struct {
	BaseHandler
}

func (h *HandlerTwo) Handle(request string) {
	if request == "Two" {
		fmt.Println("HandlerTwo обработал запрос")
	} else {
		fmt.Println("HandlerTwo не смог обработать запрос, передача дальше")
		h.BaseHandler.Handle(request)
	}
}

func main() {
	// Создаём обработчики
	handlerOne := &HandlerOne{}
	handlerTwo := &HandlerTwo{}

	// Устанавливаем цепочку вызовов
	handlerOne.SetNext(handlerTwo)

	// Передаём запросы
	fmt.Println("Передаём запрос 'One':")
	handlerOne.Handle("One")

	fmt.Println("\nПередаём запрос 'Two':")
	handlerOne.Handle("Two")

	fmt.Println("\nПередаём запрос 'Three':")
	handlerOne.Handle("Three")
}
