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
Можно воспользоваться скриптом:
```bash
chmod +x scripts/docker_run
scripts/docker_run
```

## Сборка и запуск
Установить пакеты:
```bash
go get github.com/golang/dep/cmd/dep
dep ensure
```

Сборка и запуск:
```bash
go generate -x ./restapi
go install ./cmd/forum-server
forum-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://username:password@host/database_name?sslmode=disable
```

Со скриптами:
```bash
chmod +x scripts/build
chmod +x scripts/run
./build
./run -u username -p password -u localhost(default) -d db_name
```
