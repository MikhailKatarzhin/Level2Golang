Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
ПРограмма выведет "error" по той причине, что интерфейс error, именованный как err, принимает конкретный тип *customError, и, пусть его значение остаётся nil, сам интерфейс более не сравним с nil

```
