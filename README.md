# example
go run cmd/sql-migrate/migrate.go --dir=migrations --host=percona --password=root

Выполнение sql миграций на Go во время выполнения CI тестов.

Преимущества:
* нет дополнительных зависимостей, только контейнер Go
* не нужно собирать кастомный контейнер для выполнения запросов командой Mysql
