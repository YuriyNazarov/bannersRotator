-- +goose Up
-- +goose StatementBegin
create table if not exists banners
(
    id serial
    constraint banners_pk
    primary key,
    description text
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banners;
-- +goose StatementEnd
