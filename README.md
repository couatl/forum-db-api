# Restful API to Forums
Данный репозиторий содержит реализацию API для курса СУБД в рамках ["Технопарка"](https://park.mail.ru/).

* Реализован API, описанный в [swagger.yml](https://tech-db-forum.bozaro.ru/)
* Сервис использует PostgreSQL 9.6

## Сборка и запуск Docker-контейнеров
В контейнере:

 * Приложение:
   * использует и объявляет порт 5000 (http);
 * PostgreSQL:
   * использует и объявляет порт 5342 (http);
   * директории `/etc/postgresql`, `/var/log/postgresql`, `/var/lib/postgresql` объявлены как VOLUME для возможности сохранения БД.

Сборка и запуск контейнера:
```bash
docker build -t forum-db-api -f Dockerfile .
docker run -p 5000:5000 --name forum forum-db-api
```

## Сборка и запуск
Установить пакеты:
```bash
go get github.com/golang/dep/cmd/dep
dep ensure
```

Затем для запуска скриптов:
```bash
chmod +x ./build
chmod +x ./run
./build
./run -u username -p password -u localhost(default) -d db_name
```
