package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

	Применимость: необходимо иметь ограниченный, но простой интерфейс к сложной подсистеме классов.

	Преимущества:
		- Объединяет ряд связанных последовательно или параллельно выполняемых функций в одну.
		- предоставляет только необходимый набор функций, скрывая неиспользуемые.

	Недостатки:
		- Может превратиться в антипаттерн божественный объект - излишне наполненной структурой со сложной системой, отвечающий за всё и избыточно связывая, без необходимости в том.
		- Уменьшенная функциональность всей совокупности охватываемых фасадом классов по сравнению с их прямым использованием.
*/

// Подсистемы
type Wallet struct{}

func (w *Wallet) ReduceBalance(amount float64, description string) {
	fmt.Printf("Вычтено %.2f со счёта: %s\n", amount, description)
}

type Accounting struct{}

func (a *Accounting) GeneratePackageOfDocuments(orderID string) {
	fmt.Printf("Создан пакет документов по заказу: %s\n", orderID)
}

type OrderSystem struct{}

func (o *OrderSystem) OrderPicking(orderID string) {
	fmt.Printf("Создан тикет на сборку заказа %s\n", orderID)
}

type Inventory struct{}

func (i *Inventory) UpdateStock(productID string) {
	fmt.Printf("Изменён остаток товаров согласно заказу %s\n", productID)
}

type Warehouse struct{}

func (w *Warehouse) ChangeItemsStatus(orderID string, status string) {
	fmt.Printf("Изменение статуса товаров склада согласно заказу %s на %s\n", orderID, status)
}

type Logistics struct{}

func (l *Logistics) PrepareShipment(orderID string) {
	fmt.Printf("Подготовка приказа на доставку заказа %s\n", orderID)
}

// Фасад
type OnlineStoreFacade struct {
	wallet      *Wallet
	accounting  *Accounting
	orderSystem *OrderSystem
	inventory   *Inventory
	warehouse   *Warehouse
	logistics   *Logistics
}

func NewOnlineStoreFacade() *OnlineStoreFacade {
	return &OnlineStoreFacade{
		wallet:      &Wallet{},
		accounting:  &Accounting{},
		orderSystem: &OrderSystem{},
		inventory:   &Inventory{},
		warehouse:   &Warehouse{},
		logistics:   &Logistics{},
	}
}

func (o *OnlineStoreFacade) CompletePurchaseConfirmation(orderID string, summaryCost float64) {
	fmt.Println("Processing purchase...")
	o.wallet.ReduceBalance(summaryCost, fmt.Sprintf("оплата заказа %s", orderID))
	o.accounting.GeneratePackageOfDocuments(orderID)
	o.orderSystem.OrderPicking(orderID)
	o.inventory.UpdateStock(orderID)
	o.warehouse.ChangeItemsStatus(orderID, "Reserved")
	o.logistics.PrepareShipment(orderID)
	fmt.Println("Покупка подтверждена.")
}

func main() {
	facade := NewOnlineStoreFacade()

	orderID := "ORD12345"
	summaryCost := 250.00

	facade.CompletePurchaseConfirmation(orderID, summaryCost)
}
