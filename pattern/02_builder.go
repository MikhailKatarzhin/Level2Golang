package pattern

import (
	"fmt"
	"strings"
)

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern

	Применимость:
		- Замена "телескопическому" конструктору.
		- Создание объекта из сложных вложенных структур, каждая из которых может быть создана отдельным строителем.
		- Создание отдельных представлений объекта из одинаковых этапов, которые отличаются в деталях

	Преимущества:
		- Позволяет создавать продукты пошагово.
		- Позволяет использовать один и тот же код для создания различных продуктов.
		- Изолирует сложный код сборки продукта от его основной бизнес-логики.

	Недостатки:
		- Усложняет код программы из-за введения дополнительных классов.
		- Зависимость от конкретных классов строителей.

*/

// SQLQueryBuilder — интерфейс для построения SQL-запроса.
type SQLQueryBuilder interface {
	Select(fields ...string) SQLQueryBuilder
	From(table string) SQLQueryBuilder
	Where(condition string) SQLQueryBuilder
	OrderBy(field string, asc bool) SQLQueryBuilder
	Limit(limit int) SQLQueryBuilder
	Build() string
}

// ConcreteSQLQueryBuilder — конкретный строитель SQL-запроса.
type ConcreteSQLQueryBuilder struct {
	fields     []string
	table      string
	conditions []string
	order      string
	limit      int
}

func NewSQLQueryBuilder() SQLQueryBuilder {
	return &ConcreteSQLQueryBuilder{}
}

func (b *ConcreteSQLQueryBuilder) Select(fields ...string) SQLQueryBuilder {
	b.fields = fields
	return b
}

func (b *ConcreteSQLQueryBuilder) From(table string) SQLQueryBuilder {
	b.table = table
	return b
}

func (b *ConcreteSQLQueryBuilder) Where(condition string) SQLQueryBuilder {
	b.conditions = append(b.conditions, condition)
	return b
}

func (b *ConcreteSQLQueryBuilder) OrderBy(field string, asc bool) SQLQueryBuilder {
	order := "ASC"
	if !asc {
		order = "DESC"
	}
	b.order = fmt.Sprintf("%s %s", field, order)
	return b
}

func (b *ConcreteSQLQueryBuilder) Limit(limit int) SQLQueryBuilder {
	b.limit = limit
	return b
}

func (b *ConcreteSQLQueryBuilder) Build() string {
	query := strings.Builder{}

	query.WriteString("SELECT ")
	if len(b.fields) > 0 {
		for _, field := range b.fields {
			query.WriteString(field)
			query.WriteString(", ")
		}
	} else {
		query.WriteString("*")
	}

	query.WriteString(" FROM ")
	query.WriteString(b.table)

	if len(b.conditions) > 0 {
		query.WriteString(" WHERE ")

		for _, condition := range b.conditions {
			query.WriteString(condition)
			query.WriteString(" AND ")
		}
	}
	if b.order != "" {
		query.WriteString(" ORDER BY ")
		query.WriteString(b.order)
	}
	if b.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
	}

	return query.String()
}

func main() {
	// Пример использования строителя для создания сложного SQL-запроса
	query := NewSQLQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > 18").
		Where("active = 1").
		OrderBy("name", true).
		Limit(10).
		Build()

	fmt.Println(query)
}
