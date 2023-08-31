# avitoTask
Akimov Alexandr

Cервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)


## Как поднять сервис:

```
docker-compose up --build my-app
```

## Как удалить сервис:

```
docker-compose down
```

## API. Примеры запросы:

Регистрация сегментов:
```
curl -XPOST http://0.0.0.0:8080/api/seg -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"segments": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS", "AVITO_DISCOUNT_50", "AVITO_DISCOUNT_20"]}'
```

Добавления пользователя в сегмент:
```
curl -XPOST http://0.0.0.0:8080/api/add -H 'accept: application/json' -H 'Content-Type: application/json' -d '{ "personID": 1002 , "segments": ["AVITO_VOICE_MESSAGES", "AVITO_PERFORMANCE_VAS"]}'
```

Получение сегментов у пользователя:

curl -XGET http://0.0.0.0:8080/api/person/1002 -H 'accept: application/json' -H 'Content-Type: application/json


Удаление сегмента:
```
curl -XPOST http://0.0.0.0:8080/api/deleteSegment -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"segment": ["AVITO_VOICE_MESSAGES"]}'
```

Получение пользователей определенного сегмента:
```
curl -XGET http://0.0.0.0:8080/api/segment/AVITO_VOICE_MESSAGES -H 'accept: application/json' -H 'Content-Type: application/json'
```

Удалить сегменты у пользователя:

```
curl -XPOST http://0.0.0.0:8080/api/deleteSegments/1002 -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"segments": ["AVITO_PERFORMANCE_VAS"]}'
```


## Как работает сервис


