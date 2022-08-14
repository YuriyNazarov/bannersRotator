-- +goose Up
-- +goose StatementBegin
create table if not exists dem_groups
(
    id serial
    constraint dem_groups_pk
    primary key,
    description text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table dem_groups;
-- +goose StatementEnd
