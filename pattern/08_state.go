package pattern

import (
	"fmt"
	"time"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

	Применимость:
		- Есть объект, поведение которого кардинально меняется в зависимости от внутреннего состояния, причём типов состояний много, и их код часто меняется.
		- Когда код класса содержит множество больших, похожих друг на друга, условных операторов, которые выбирают поведения в зависимости от текущих значений полей класса.
		- Когда вы сознательно используете табличную машину состояний, построенную на условных операторах, но вынуждены мириться с дублированием кода для похожих состояний и переходов.

	Преимущества:
		- Избавляет от множества больших условных операторов машины состояний.
		- Концентрирует в одном месте код, связанный с определённым состоянием.
		- Упрощает код контекста.

	Недостатки:
		- Может неоправданно усложнить код, если состояний мало и они редко меняются.
		- Трудность управления сложными переходами между состояниями.

*/

// TrafficLightState Интерфейс состояния
type TrafficLightState interface {
	ChangeLight(context *TrafficLight)
	ShowLight()
}

// RedLight Конкретное состояние Красный свет
type RedLight struct{}

func (r *RedLight) ChangeLight(context *TrafficLight) {
	fmt.Println("Меняем состояние с Красного на Желтый (перед зелёным).")
	context.SetState(&YellowBeforeGreenLight{})
}

func (r *RedLight) ShowLight() {
	fmt.Println("Красный свет. Стоп.")
}

// YellowBeforeGreenLight Конкретное состояние  Желтый свет перед зелёным
type YellowBeforeGreenLight struct{}

func (y *YellowBeforeGreenLight) ChangeLight(context *TrafficLight) {
	fmt.Println("Меняем состояние с Желтого на Зеленый.")
	context.SetState(&GreenLight{})
}

func (y *YellowBeforeGreenLight) ShowLight() {
	fmt.Println("Желтый свет. Приготовьтесь.")
}

// GreenLight Конкретное состояние Зеленый свет
type GreenLight struct{}

func (g *GreenLight) ChangeLight(context *TrafficLight) {
	fmt.Println("Меняем состояние с Зеленого на Желтый.")
	context.SetState(&YellowBeforeRedLight{})
}

func (g *GreenLight) ShowLight() {
	fmt.Println("Зеленый свет. Можно ехать.")
}

// YellowBeforeRedLight Конкретное состояние Желтый свет перед красным
type YellowBeforeRedLight struct{}

func (y *YellowBeforeRedLight) ChangeLight(context *TrafficLight) {
	fmt.Println("Меняем состояние с Желтого на Красный.")
	context.SetState(&RedLight{})
}

func (y *YellowBeforeRedLight) ShowLight() {
	fmt.Println("Желтый свет. Скоро будет красный.")
}

// TrafficLight Контекст, содержащий текущее состояние
type TrafficLight struct {
	state TrafficLightState
}

func (t *TrafficLight) SetState(state TrafficLightState) {
	t.state = state
}

func (t *TrafficLight) ChangeLight() {
	t.state.ChangeLight(t)
}

func (t *TrafficLight) ShowLight() {
	t.state.ShowLight()
}

func main() {
	// Начальное состояние светофора — красный свет
	trafficLight := &TrafficLight{state: &RedLight{}}

	for i := 0; i < 8; i++ { // Цикл для демонстрации работы светофора
		trafficLight.ShowLight()    // Показываем текущий свет
		time.Sleep(2 * time.Second) // Имитация ожидания
		trafficLight.ChangeLight()  // Меняем состояние
	}
}
