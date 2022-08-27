-- +goose Up
-- +goose StatementBegin
create table if not exists actions
(
    id serial
        constraint actions_pk
            primary key,
    action_type int not null,
    slot_id int not null,
    dem_group_id int not null,
    created_at timestamp default current_timestamp not null,
    banner_id int not null
);

comment on column actions.action_type is '0 - view
1 - click';


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table actions;
-- +goose StatementEnd
