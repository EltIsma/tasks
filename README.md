# Tasks
## Описание


## API

```

```

## Установка и запуск

```
make compose

# Ждём пока все запустится... 
# Отслеживать можно по логам

# Применяем миграции к Postgres:
migrate -path migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'  up

# Создаем топик:
make create-kafka-topic

=> Profit!
```