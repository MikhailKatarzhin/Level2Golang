Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
Программа выведет построчно <nil> и false, причиной чего является то, что 
1) интерфейс, чем является error, состоит из конкретного типа и указателя на значение;
2) fmt.Println(err) выведет значение интерфейса, т.е. указатель *os.PathError
3) fmt.Println(err == nil) выведет Ложь, так как конкретный тип указан *os.PathError


```
