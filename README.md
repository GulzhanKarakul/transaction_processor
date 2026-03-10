# Transaction Processor

HTTP сервис параллельной обработки финансовых транзакций.

## Стек

Go · Worker Pool · Pipeline · sync.RWMutex · Context · Docker

## Архитектура
```
HTTP запрос → Handler → Worker Pool (5 горутин) → Pipeline → Storage
                                                   ↓
                                        validate → addBonus → save
```

## Запуск
```bash
docker-compose up --build
```

## API

### Создать транзакцию
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-1", "amount": 1000}'
```

### Получить все транзакции
```bash
curl http://localhost:8080/transactions
```

## Тесты
```bash
go test -race ./...
```