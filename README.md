# avitoTask
Автор: Акимов Александр

Cервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)

Также написан [Swagger](https://github.com/brokensm1le/avitoTask/blob/main/spec.yaml)


## Работа с сервисом:

#### Как поднять сервис:

```
docker-compose up --build -d my-app
```

#### Kак подключиться к БД вручную:

```
sudo apt-get install postgresql-client
psql -h 127.0.0.1 -U root -d taskdb
```

#### Как удалить сервис:

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
```
curl -XGET http://0.0.0.0:8080/api/person/1002 -H 'accept: application/json' -H 'Content-Type: application/json
```

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

### Доп задание №1

```
curl -XGET http://0.0.0.0:8080/api/checkHistory/1002 -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"timeFrom": "2023-08-31T18:05:12.000Z", "timeTo": "2023-08-31T18:13:40.000Z"}'
```

*Замечание.* По тексту задания не понял какую именно ссылку надо было отправлять пользователю. Поэтому вывод истории выводится в консоль.


### Доп задание №3

В данном задание просто рандомил число от 1 до 100, если число больше чем *percentage*, то сегмент не попал к Пользователю, иначе попал.

```
curl -XPOST http://0.0.0.0:8080/api/addWithPercentage -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"segment": "AVITO_PERFORMANCE_VAS", "IDs": [1003, 1004, 1005, 1006], "percentage":50}'
```

### Доп задание №2

*(Не выполнил)* Но идея выполнения состоит в том чтобы создать сервис, который отправляет запросы к БД через определенные промежутки времени(условно, каждые 2 минуты). И смотрит можно ли что-то удалить.

## Подробности о сервисе

### БД
В данном проекте используется БД Postrgesql;

Всего две таблицы:  
 - id_seg: которая хранит информацию id пользоваетля -> сегменты пользоваля
 - seg_id: которая хранит информацию сегмент -> id пользователей у которых есть данный сегмент
 - hystiry: для хранения истории запросов(доп задание №1)

Это было сделано для того, чтобы уменьшить количество запросов. Те при удалении сегмента нам будет необходимо пройтись не по всей Таблице а только по определенным ID.

Так как мы хотим чтобы при "условном" падении нашего сервиса или БД, в базе данных осталась верная информация используются транзакции.
