# Worktime
## Консольная утилита для учета рабочего времени

### Описание
```
Использование: worktime (start|stop|time [full]|dinner (minutes))
   start            Отметка о начале рабочего дня
   stop             Отметка об окончании рабочего дня
   dinner (minutes) Запись количества минут проведенных на отдыхе или обеде
   time             Просмотр временного баланса переработок или недоработок
   time full        Просморт полного лога рабочего времени
   help             Просмотр текущей справки
```

### Установка
* Установите [Golang](https://golang.org/doc/install)
* Склонируйте этот репозиторй
* Скачайте необходимые зависимости `go get ./...`
* Выполните команду `go build`
* Исполняемый файл ./worktime готов к запуску