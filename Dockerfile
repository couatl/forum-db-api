FROM ubuntu:16.04

RUN apt-get update

ENV PGVER 9.5
RUN apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

RUN apt install -y golang-1.8 git

ENV GOROOT /usr/lib/go-1.8
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/github.com/couatl/forum-db-api
ADD . $GOPATH/src/github.com/couatl/forum-db-api

RUN go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
RUN go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata

# Собираем и устанавливаем пакет
RUN go generate -x ./restapi
RUN go install ./cmd/forum-server

EXPOSE 5000

CMD service postgresql start && forum-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker
