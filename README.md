# Restful API to Forums
Данный репозиторий содержит реализацию API для курса СУБД в рамках ["Технопарка"](https://park.mail.ru/).

* Реализован API, описанный в [swagger.yml](https://tech-db-forum.bozaro.ru/)
* Сервис использует PostgreSQL 9.6

## Сборка и запуск Docker-контейнеров
Docker контейнера организованы по следующему приципу:

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

Остановить контейнер можно командой:
```bash
docker kill forum
```

## Сборка и запуск
Установить пакеты:
```bash
go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata
```

Затем:
```bash
go generate -x ./restapi
go install ./cmd/forum-server
forum-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://username:password@localhost/db_name?sslmode=disable
```
