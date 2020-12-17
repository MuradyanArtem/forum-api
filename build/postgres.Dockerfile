FROM postgres:alpine

COPY assets/dump.sql /dump.sql

COPY scripts/initdb.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/create-db.sh

EXPOSE 5432
