# [Тестовое задание](https://github.com/avito-tech/mx-backend-trainee-assignment)

## Запуск

```sh
git clone https://github.com/adzhukov/avito_mx.git
cd avito_mx
docker-compose up
```

## API

### /import/sync

Синхронный импорт:

- `url` – URI .xlsx файла с товарами
- `seller_id` – ID продавца

###### Запрос

```sh
curl http://localhost:8080/import/sync -G \
  -d seller_id=11 \
  -d url=file:///files/10k.xlsx
```

###### Ответ

```json
{
  "task_id": 1875,
  "status": "Success",
  "stats": {
    "created": 0,
    "updated": 4999,
    "deleted": 0,
    "invalid": 1
  }
}
```

### /import

Асинхронный импорт:

- `url` – URI .xlsx файла с товарами
- `seller_id` – ID продавца

###### Запрос

```sh
curl http://localhost:8080/import -G \
  -d seller_id=11 \
  -d url=file:///files/10k.xlsx
```

###### Ответ

```json
{
  "task_id": 1873,
  "status": "Queued"
}
```

### /status

Проверка статуса задачи

Параметры:

- `task_id` – ID задачи из ответа `/import`

###### Запрос

```sh
curl http://localhost:8080/status -G \
  -d task_id=1873
```

###### Ответ

```json
{
  "task_id": 1873,
  "status": "Success",
  "stats": {
    "created": 0,
    "updated": 4999,
    "deleted": 0,
    "invalid": 1
  }
}
```

### /offers

Получение списка товаров из базы

Параметры:

- `seller_id` – ID продавца
- `offer_id` – ID товара в системе продавца
- `q` – Подстрока поиска

###### Запрос

```sh
curl http://localhost:8080/offers -G \
  -d seller_id=11 \
  -d offer_id=8 \
  -d q=Offe
```

###### Ответ

```json
[
  {
    "seller_id": 11,
    "offer_id": 8,
    "name": "Offer 00008",
    "price": 8000,
    "quantity": 800
  }
]
```
