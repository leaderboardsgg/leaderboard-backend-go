FROM postgres:alpine3.13

ENV POSTGRES_PASSWORD=password

COPY schema.sql /schema.sql
COPY testdata.sql /testdata.sql
COPY init.sh /docker-entrypoint-initdb.d/init.sh
