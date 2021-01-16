FROM postgres:12
ENV POSTGRES_DB clinicdb
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
COPY init.sql /docker-entrypoint-initdb.d/
COPY 1_initiali_tabels.tx.up.sql /docker-entrypoint-initdb.d/
COPY 2_mock.sql /docker-entrypoint-initdb.d/