# Tinder
Аналог популярного сервиса для онлайн-знакомств.

## Описание
Микросервисная масшабируемая архитектура
Все микросервисы Stateless
Асинхронное и синхронное виды взаимодействия микросервисов
PROD качество: observability (logging, metrics, tracing), graceful shutdown, testing (unit, integration, e2e)
Разделенные потоки чтения и изменения данных (CQRS)
Стек: Golang, C#, Kafka, Rest, Grpc, PostgreSQL, Redis, Docker, Docker compose, Jaeger
Архитектурные подходы: outbox, eventual consistency, stateless, distributed tracing
JWT авторизация
Docker compose для запуска системы и всех ее зависимостей: /cicd/local/docker-compose.yml

На данный момент реализованы ProfileService на C# и GeoService на Golang.
ProfileService отвечает за хранение и управление профилями пользователей.
GeoService использует подход [геохеширования](https://en.wikipedia.org/wiki/Geohash) для быстрого поиска ближайших пользователей.
В процессе реализация ReactionService и MatchService на Golang.

## System Design
![alt text](https://github.com/nsinitsyn/tndr/blob/master/architecture/system%20design.png?raw=true)