-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tg_users
(
    tg_id    INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS cart
(
    id         INTEGER PRIMARY KEY GENERATED ALWAYS AS identity ( increment by 1 start 1),
    owner_id   INTEGER REFERENCES tg_users (tg_id)
    );

CREATE TABLE IF NOT EXISTS carts_items
(
    cart_id        INTEGER PRIMARY KEY GENERATED ALWAYS AS identity ( increment by 1 start 1),
    item_name      text[],
    user_id        INTEGER REFERENCES tg_users (tg_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carts_items;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS tg_users;
-- +goose StatementEnd
