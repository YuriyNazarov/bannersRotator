-- +goose Up
-- +goose StatementBegin
create table if not exists slots
(
    id serial
        constraint slots_pk
            primary key,
    description text
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table slots;
-- +goose StatementEnd
