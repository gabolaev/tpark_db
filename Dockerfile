FROM ubuntu:17.10
MAINTAINER George Gabolaev

ENV PGVERSION=9.6
ENV GOVERSION=1.10
ENV REPO=github.com/gabolaev/tpark_db

# Basic tools
RUN apt update
RUN apt install -y git vim wget

# PostgreSQL
RUN apt install -y postgresql-$PGVERSION postgresql-contrib

# Create user and database
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    /etc/init.d/postgresql stop

USER root
# Open Postgres for network
RUN echo "local all all md5" > /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVERSION/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# config
ENV PGHOST /var/run/postgresql
ENV PGDATABASE forum
ENV PGUSER forum
ENV PGPASSWORD forum
EXPOSE 5432

# GoLang
RUN wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg
ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
ADD ./ $GOPATH/src/$REPO
WORKDIR $GOPATH/src/$REPO
RUN go install .
RUN go build
EXPOSE 5000

RUN echo "./config/postgresql.conf" >> /etc/postgresql/$PGVERSION/main/postgresql.conf
USER postgres
CMD /etc/init.d/postgresql start && ./tpark_db


