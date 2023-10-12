create schema inventory;
create table inventory.items(
    id  SERIAL PRIMARY KEY,
    name        TEXT    NOT NULL,
    description   TEXT    NOT NULL,
    price       DECIMAL(10,2) NOT NULL, 
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);