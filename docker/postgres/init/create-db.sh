#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
		CREATE USER wb WITH PASSWORD 'wb';
		CREATE DATABASE wb;
		GRANT ALL ON SCHEMA public TO wb;
		GRANT ALL ON DATABASE wb TO wb;
		ALTER DATABASE wb OWNER TO wb;
EOSQL

psql -U wb -d wb -f /init.sql
