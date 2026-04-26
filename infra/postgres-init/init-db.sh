#!/bin/bash
set -e

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		echo "  Creating database '$db'"
		psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
		    CREATE DATABASE $db;
EOSQL
	done
	echo "Multiple databases created"
	
	echo "Initializing payment_db schema..."
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "payment_db" <<-EOSQL
		CREATE TABLE IF NOT EXISTS payments (
			id UUID PRIMARY KEY,
			order_id UUID NOT NULL,
			amount DECIMAL NOT NULL,
			status VARCHAR(20) NOT NULL,
			idempotency_key VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
EOSQL

	echo "Initializing order_db schema (Outbox)..."
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "order_db" <<-EOSQL
		CREATE TABLE IF NOT EXISTS outbox_events (
			id UUID PRIMARY KEY,
			aggregate_id UUID NOT NULL,
			aggregate_type VARCHAR(50) NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			payload TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
EOSQL
fi
