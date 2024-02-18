-- +goose Up
-- +goose StatementBegin
ALTER TABLE carts_items add column status varchar(4)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
