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

Docker compose для запуска системы и всех ее зависимостей: /cicd/local/docker-compose.yml

На данный момент реализованы ProfileService на C# и GeoService на Golang.
ProfileService отвечает за хранение и управление профилями пользователей.
GeoService использует подход [геохеширования](https://en.wikipedia.org/wiki/Geohash) для быстрого поиска ближайших пользователей.
В процессе реализация ReactionService и MatchService на Golang.

## System Design
![alt text](https://github.com/nsinitsyn/tndr/blob/master/architecture/system%20design.png?raw=true)

## Скриншоты
Пример обработки конфликта параллелизма при оптимистичном обновлении данных в Redis - скриншот трейса из Jaeger, полученный на end-to-end тесте
![alt text](https://github.com/nsinitsyn/tndr/blob/master/architecture/redis%20optimistic%20locking%20-%20jeager.png?raw=true)