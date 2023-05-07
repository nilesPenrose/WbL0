FROM postgres:latest

COPY docker/postgres/init/init.sql /
COPY docker/postgres/init/create-db.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/create-db.sh

EXPOSE 5432
