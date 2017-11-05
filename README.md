# Restful API to Forums
Данный репозиторий содержит реализацию API для курса СУБД в рамках ["Технопарка"](https://park.mail.ru/pages/index/).

* Реализован API, описанный в swagger.yml
* Сервис использует PostgreSQL 9.5

## Состав Docker-контейнеров
Docker контейнера организованы по следующему приципу:

 * Приложение:
   * использует и объявляет порт 5000 (http);
 * PostgreSQL:
   * использует и объявляет порт 5342 (http);
   * директории `/etc/postgresql`, `/var/log/postgresql`, `/var/lib/postgresql` объявлены как VOLUME для возможности сохранения БД.

## Как собрать и запустить контейнер

Сборка и запуск контейнера:
```bash
docker build -t forum-db-api -f Dockerfile .
docker run -p 5000:5000 --name forum forum-db-api
```

Остановить контейнер можно командой:
```bash
docker kill forum
```
