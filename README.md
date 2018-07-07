# Worktime
## Консольная утилита для учета рабочего времени

[![Go Report Card](https://goreportcard.com/badge/github.com/Vehsamrak/worktime)](https://goreportcard.com/report/github.com/Vehsamrak/worktime)

### Описание
```
Использование: worktime (start|stop|time [full]|dinner (minutes)|note (comment)|version)
   start 		Отметка о начале рабочего дня
   stop 		Отметка об окончании рабочего дня
   dinner (minutes) 	Запись количества минут проведенных на отдыхе или обеде
   time 		Просмотр временного баланса переработок или недоработок
   time full		Просморт полного лога рабочего времени
   time full [X]	Просморт лога рабочего времени за X последних дней
   note (text comment) 	Добавление комментария к текущему дню
   version 		Отображение текущей версии
   help 		Просмотр текущей справки
```

### Установка
Скомпиллированное и готовое к запуску приложение доступно в [последнем релизе](https://github.com/Vehsamrak/worktime/releases/latest).

### Сборка из исходников
* Установите [Golang](https://golang.org/doc/install)
* Склонируйте этот репозиторй
* Скачайте необходимые зависимости `go get ./...`
* Выполните команду `go build`
* Исполняемый файл ./worktime готов к запуску
