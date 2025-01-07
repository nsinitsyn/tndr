# Tinder
Аналог популярного сервиса для онлайн-знакомств.

## Описание
Микросервисная масштабируемая архитектура

Все микросервисы Stateless

Асинхронное и синхронное виды взаимодействия микросервисов

PROD качество: observability (logging, metrics, tracing), graceful shutdown, testing (unit, integration, e2e)

Разделенные потоки чтения и изменения данных (CQRS)

Стек: Golang, C#, Kafka, Rest, Grpc, PostgreSQL, Redis, Docker, Docker compose, Jaeger

Архитектурные подходы: outbox, eventual consistency, stateless, distributed tracing, optimistic locking

JWT авторизация

Docker compose для запуска системы и всех ее зависимостей: `cicd/local/docker-compose.yml`

На данный момент реализованы ProfileService на C# и GeoService на Golang.
ProfileService отвечает за хранение и управление профилями пользователей.
GeoService использует подход [геохеширования](https://en.wikipedia.org/wiki/Geohash) для быстрого поиска ближайших пользователей.
В процессе реализация ReactionService и MatchService на Golang.

## System Design
![alt text](https://github.com/nsinitsyn/tndr/blob/master/architecture/system%20design.png?raw=true)

## Скриншоты
Пример обработки конфликта параллелизма при оптимистичном обновлении данных в Redis - скриншот трейса из Jaeger, полученный на end-to-end тесте
![alt text](https://github.com/nsinitsyn/tndr/blob/master/architecture/redis%20optimistic%20locking%20-%20jeager.png?raw=true)

## Запуск приложения
### Docker
Geo service с нужной инфраструктурой запускается через docker compose `cicd/local/docker-compose.yml`

После запуска он будет слушать входящие grpc запросы на порту `2342`. Метрики доступны по HTTP на `2322/metrics`

Для тестирования можно отправить следующий grpc-запрос через grpcurl (тестовый jwt-токен заранее сгенерирован на длительный срок):
```
grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxIiwiR2VuZGVyIjoiTSIsImV4cCI6MTc2MzIwNzQxMywiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.VAVP65lIUhabxR4UknvQkRKiVCfu116cf3tZC8-dsfw' -plaintext -d '{"latitude":55.481, "longitude":37.288}' localhost:2342 tinder.GeoService/GetProfilesByLocation
```

Должен быть получен пустой ответ, т.к. не были созданы профили (эта функция тоже работает)

В jaeger можно увидеть трассировку по данному запросу с командой в Redis
### MAC OS
Для запуска инфраструктуры и приложения необходимо использовать следующие команды в терминале:
```
# run infra
cd cicd/debug
docker compose up -d
# run application
cd src/Geo
export CONFIG_PATH='config/config.yaml'; export CGO_ENABLED=1; go run cmd/tinder-geo/main.go
```
### Windows
Предварительно должен быть установлен GCC. В данном примере используется [TDM-GCC](https://github.com/jmeubank/tdm-gcc), который установлен по пути `C:\TDM-GCC-64\bin\gcc`.
Для запуска инфраструктуры и приложения необходимо использовать следующие команды в командной строке:
```
# run infra
cd cicd/debug
docker compose up -d
# run application
cd src/Geo
$env:CONFIG_PATH = '././config/config.yaml'; $env:CGO_ENABLED=1; $env:CC="C:\TDM-GCC-64\bin\gcc"; go run cmd/tinder-geo/main.go
```