package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

	Паттерн «Посетитель» часто используется в таких задачах, как обход деревьев выражений, где к узлам применяется одна и та же операция. Также его используют при обработке документов или в графических редакторах, например, когда нужно добавлять или изменять операции над разными типами элементов.

	Применимость:
		- Когда нужно выполнить операцию над всеми элементами сложной структуры объектов, например, деревом.
		- Когда новое поведение имеет смысл только для некоторых классов из существующей иерархии, позволяя определить поведение только для этих классов, оставив его пустым для всех остальных

	Преимущества:
		- Упрощает добавление операций, работающих со сложными структурами объектов.
		- Объединяет родственные операции в одном классе.
		- Посетитель может накапливать состояние при обходе структуры элементов.
		- Методы можно вынести из классов структуры в отдельные классы посетителей, что упрощает сами классы.

	Недостатки:
		- Паттерн не оправдан, если иерархия элементов часто меняется. так как потребуется обновить всех посетителей.

*/

// Shape Интерфейс для всех фигур
type Shape interface {
	Accept(visitor ShapeVisitor)
}

// Circle Конкретный класс интерфейса Shape
type Circle struct {
	Radius float64
}

func (c *Circle) Accept(visitor ShapeVisitor) {
	visitor.VisitCircle(c)
}

// Rectangle Конкретный класс интерфейса Shape
type Rectangle struct {
	Width, Height float64
}

func (r *Rectangle) Accept(visitor ShapeVisitor) {
	visitor.VisitRectangle(r)
}

// ShapeVisitor Интерфейс посетителя
type ShapeVisitor interface {
	VisitCircle(*Circle)
	VisitRectangle(*Rectangle)
}

// AreaCalculator Конкретный посетитель для вычисления площади
type AreaCalculator struct {
	Area float64
}

func (a *AreaCalculator) VisitCircle(c *Circle) {
	a.Area = 3.14 * c.Radius * c.Radius
	fmt.Printf("Площадь круга: %.2f\n", a.Area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) {
	a.Area = r.Width * r.Height
	fmt.Printf("Площадь прямоугольника: %.2f\n", a.Area)
}

func main() {
	// Массив фигур
	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 3, Height: 4},
		&Circle{Radius: 2.5},
		&Rectangle{Width: 5, Height: 6},
	}

	areaCalculator := &AreaCalculator{}

	// Проходим по каждой фигуре и вызываем Accept для каждой из них
	for _, shape := range shapes {
		shape.Accept(areaCalculator)
	}
}
