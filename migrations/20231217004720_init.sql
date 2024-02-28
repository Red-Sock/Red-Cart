-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tg_users (
    tg_id      INTEGER PRIMARY KEY,
    user_name  TEXT,
    first_name TEXT,
    last_name  TEXT
);

CREATE TABLE IF NOT EXISTS carts (
    id         INTEGER PRIMARY KEY GENERATED ALWAYS AS identity ( increment by 1 start 1),
    owner_id   INTEGER REFERENCES tg_users (tg_id),
    chat_id    INTEGER,
    message_ID INTEGER
);

CREATE TABLE IF NOT EXISTS cart_items (
    cart_id   INTEGER REFERENCES carts (id),
    item_name TEXT,
    amount    INTEGER,
    user_id   INTEGER REFERENCES tg_users (tg_id),
    UNIQUE (user_id, cart_id, item_name)
);

CREATE TABLE IF NOT EXISTS carts_users (
    user_id    INTEGER REFERENCES tg_users (tg_id),
    cart_id    INTEGER REFERENCES carts (id),
    is_default BOOLEAN
);

CREATE UNIQUE INDEX ON carts_users (is_default) WHERE is_default = TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS tg_users;
-- +goose StatementEnd
