-- +goose Up
Alter Table feeds Add Column last_fetched_at TIMESTAMP NULL DEFAULT NULL;

-- +goose Down
Alter Table feeds Drop Column last_fetched_at;

