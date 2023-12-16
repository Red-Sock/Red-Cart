-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tg_users
(
    tg_id    INTEGER PRIMARY KEY,
    user_name varchar,
    first_name varchar,
    last_name varchar
);

CREATE TABLE IF NOT EXISTS cart
(
    id         INTEGER PRIMARY KEY GENERATED ALWAYS AS identity ( increment by 1 start 1),
    owner_id   INTEGER REFERENCES tg_users (tg_id)
    );

CREATE TABLE IF NOT EXISTS carts_items
(
    cart_id        INTEGER REFERENCES cart (id),
    item_name      text[],
    user_id        INTEGER REFERENCES tg_users (tg_id),
    UNIQUE (user_id, cart_id)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tg_users;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS carts_items;
-- +goose StatementEnd
