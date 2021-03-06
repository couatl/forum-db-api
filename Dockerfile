FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y wget git

ENV PGVER 9.6

RUN apt-get install -y software-properties-common python-software-properties
RUN add-apt-repository "deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main"
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update
RUN apt-get install -y postgresql-$PGVER

# RUN apt-get install -y postgresql-$PGVER
# RUN apt-get install -y golang git

RUN wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz && mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql -c "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "autovacuum = off" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "work_mem = 64MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/vaer/log/postgresql", "/var/lib/postgresql"]

USER root

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/github.com/couatl/forum-db-api
ADD . $GOPATH/src/github.com/couatl/forum-db-api

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
RUN go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata

# Собираем и устанавливаем пакет
RUN go generate -x ./restapi
RUN go install ./cmd/forum-server

EXPOSE 5000

CMD service postgresql start && forum-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker
