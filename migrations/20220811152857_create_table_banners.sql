-- +goose Up
-- +goose StatementBegin
create table if not exists banners
(
    id serial
    constraint banners_pk
    primary key,
    title varchar(255) not null,
    description text
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banners;
-- +goose StatementEnd
