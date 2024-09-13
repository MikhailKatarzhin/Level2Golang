package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

	Применимость:
		- История действий и их отмена.
		- Параметризация объектов командами.
		- Последовательное выполнение команд, которые могут быть отменены в случае ошибки.

	Преимущества:
		- Убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их непосредственно выполняют.
 		- Позволяет реализовать простую отмену и повтор операций.
		-  Позволяет реализовать отложенный запуск операций.
		-  Позволяет собирать сложные команды из простых.

	Недостатки:
		- В случаях, когда требуется небольшое количество простых команд, может возникнуть избыточность.
		- Каждый новый тип команды требует создания отдельного класса, что может усложнять систему.


*/

// Command Интерфейс команды
type Command interface {
	Execute()
	Undo()
}

// Light Получатель команды
type Light struct {
	isOn       bool
	brightness int // Яркость света в процентах (от 0 до 100)
}

func (l *Light) On() {
	l.isOn = true
	fmt.Println("Свет включен на яркости", l.brightness, "%")
}

func (l *Light) Off() {
	l.isOn = false
	fmt.Println("Свет выключен")
}

func (l *Light) SetBrightness(brightness int) {
	l.brightness = brightness
	if l.isOn {
		fmt.Println("Яркость установлена на", l.brightness, "%")
	}
}

// TurnOnCommand Конкретные команды для включения света
type TurnOnCommand struct {
	light          *Light
	prevIsOnState  bool
	prevBrightness int
}

func (c *TurnOnCommand) Execute() {
	// Сохраняем предыдущее состояние и яркость
	c.prevIsOnState = c.light.isOn
	c.prevBrightness = c.light.brightness
	c.light.On()
}

func (c *TurnOnCommand) Undo() {
	// Восстанавливаем предыдущее состояние и яркость
	c.light.brightness = c.prevBrightness
	if c.prevIsOnState {
		c.light.On()
	} else {
		c.light.Off()
	}
}

// TurnOffCommand Конкретные команды для выключения света
type TurnOffCommand struct {
	light          *Light
	prevIsOnState  bool
	prevBrightness int
}

func (c *TurnOffCommand) Execute() {
	// Сохраняем предыдущее состояние и яркость
	c.prevIsOnState = c.light.isOn
	c.prevBrightness = c.light.brightness
	c.light.Off()
}

func (c *TurnOffCommand) Undo() {
	// Восстанавливаем предыдущее состояние и яркость
	c.light.brightness = c.prevBrightness
	if c.prevIsOnState {
		c.light.On()
	} else {
		c.light.Off()
	}
}

// SetBrightnessCommand Команда для изменения яркости
type SetBrightnessCommand struct {
	light          *Light
	prevBrightness int
	newBrightness  int
}

func (c *SetBrightnessCommand) Execute() {
	// Сохраняем текущее значение яркости
	c.prevBrightness = c.light.brightness
	c.light.SetBrightness(c.newBrightness)
}

func (c *SetBrightnessCommand) Undo() {
	// Восстанавливаем предыдущее значение яркости
	c.light.SetBrightness(c.prevBrightness)
}

// RemoteControl Инициатор, который вызывает команды
type RemoteControl struct {
	buttons      map[string]Command
	emergencyOff Command
}

func NewRemoteControl() *RemoteControl {
	return &RemoteControl{buttons: make(map[string]Command)}
}

func (r *RemoteControl) SetCommand(buttonName string, command Command) {
	r.buttons[buttonName] = command
}

func (r *RemoteControl) PressButton(buttonName string) {
	if command, ok := r.buttons[buttonName]; ok {
		command.Execute()
	} else {
		fmt.Println("Кнопка не назначена:", buttonName)
	}
}

func (r *RemoteControl) PressUndo(buttonName string) {
	if command, ok := r.buttons[buttonName]; ok {
		command.Undo()
	} else {
		fmt.Println("Кнопка не назначена:", buttonName)
	}
}

func (r *RemoteControl) PressEmergencyOff() {
	r.emergencyOff.Execute()
}

func main() {
	light := &Light{brightness: 50} // Начальная яркость — 50%

	// Создаем команды
	turnOn := &TurnOnCommand{light: light}
	turnOff := &TurnOffCommand{light: light}
	setBrightnessHigh := &SetBrightnessCommand{light: light, newBrightness: 100}
	setBrightnessLow := &SetBrightnessCommand{light: light, newBrightness: 10}

	// Создаем пульт с кнопками
	remote := NewRemoteControl()
	remote.SetCommand("Power On", turnOn)
	remote.SetCommand("Power Off", turnOff)
	remote.SetCommand("Bright 100%", setBrightnessHigh)
	remote.SetCommand("Bright 10%", setBrightnessLow)

	// Назначаем кнопку экстренного отключения (можно выключить свет независимо от текущего состояния)
	remote.emergencyOff = turnOff

	// Используем кнопки
	fmt.Println("Нажимаем кнопку 'Power On':")
	remote.PressButton("Power On")

	fmt.Println("\nУстанавливаем яркость на 100%:")
	remote.PressButton("Bright 100%")

	fmt.Println("\nУстанавливаем яркость на 10%:")
	remote.PressButton("Bright 10%")

	fmt.Println("\nОтмена последнего действия (возврат к предыдущей яркости):")
	remote.PressUndo("Bright 10%")

	fmt.Println("\nНажимаем кнопку экстренного отключения света:")
	remote.PressEmergencyOff()

	fmt.Println("\nОтмена экстренного отключения:")
	remote.PressUndo("Power Off")
}
