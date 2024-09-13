package pattern

import (
	"fmt"
	"time"
)

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

	Применимость:
		- Когда у вас есть множество похожих классов, отличающихся только некоторым поведением.
		- Когда вы не хотите обнажать детали реализации алгоритмов для других классов.
		- Когда различные вариации алгоритмов реализованы в виде развесистого условного оператора. Каждая ветка такого оператора представляет собой вариацию алгоритма.
		- Когда вам нужно использовать разные вариации какого-то алгоритма внутри одного объекта.

	Преимущества:
		- Горячая замена алгоритмов на лету.
 		- Изолирует код и данные алгоритмов от остальных классов.
 		- Уход от наследования к делегированию.

	Недостатки:
		- Усложняет программу за счёт дополнительных классов.
 		- Клиент должен знать, в чём состоит разница между стратегиями, чтобы выбрать подходящую.

*/

// RouteStrategy - интерфейс для различных стратегий построения маршрута
type RouteStrategy interface {
	BuildRoute(currentLocation, destination string) string
}

// ShortestRoute - стратегия для нахождения кратчайшего пути
type ShortestRoute struct{}

func (s *ShortestRoute) BuildRoute(currentLocation, destination string) string {
	// Реализация алгоритма для поиска кратчайшего маршрута
	return fmt.Sprintf("Building shortest route from %s to %s", currentLocation, destination)
}

// FastestRoute - стратегия для нахождения самого быстрого пути
type FastestRoute struct{}

func (f *FastestRoute) BuildRoute(currentLocation, destination string) string {
	// Реализация алгоритма для поиска самого быстрого маршрута с учетом пробок
	return fmt.Sprintf("Building fastest route from %s to %s (considering traffic)", currentLocation, destination)
}

// CheapestRoute - стратегия для нахождения самого дешевого маршрута
type CheapestRoute struct{}

func (c *CheapestRoute) BuildRoute(currentLocation, destination string) string {
	// Реализация алгоритма для поиска самого дешевого маршрута, избегая платных дорог
	return fmt.Sprintf("Building cheapest route from %s to %s (minimizing toll roads)", currentLocation, destination)
}

// Navigator - контекст, который использует стратегию для построения маршрута
type Navigator struct {
	strategy        RouteStrategy
	currentLocation string
}

func (n *Navigator) SetStrategy(strategy RouteStrategy) {
	fmt.Println("Switching routing strategy...")
	n.strategy = strategy
}

func (n *Navigator) SetCurrentLocation(location string) {
	n.currentLocation = location
}

func (n *Navigator) BuildRoute(destination string) {
	// Передаем текущую позицию в стратегию
	fmt.Println(n.strategy.BuildRoute(n.currentLocation, destination))
}

func main() {
	// Начальное текущее положение
	currentLocation := "Current Location"

	// Создаем навигатор с начальной стратегией
	navigator := &Navigator{
		currentLocation: currentLocation,
	}

	// Устанавливаем стратегию кратчайшего маршрута
	navigator.SetStrategy(&ShortestRoute{})
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(2 * time.Second)
			navigator.BuildRoute("City B")
		}
	}()

	// Через 3 секунды переключимся на самый быстрый маршрут
	time.Sleep(3 * time.Second)
	navigator.SetStrategy(&FastestRoute{})
	go func() {
		for i := 0; i < 2; i++ {
			time.Sleep(2 * time.Second)
			navigator.BuildRoute("City B")
		}
	}()

	// Через 6 секунд переключимся на самый дешевый маршрут
	time.Sleep(6 * time.Second)
	navigator.SetStrategy(&CheapestRoute{})
	navigator.BuildRoute("City B")

	time.Sleep(5 * time.Second) // даем системе время для выполнения всех операций
}
