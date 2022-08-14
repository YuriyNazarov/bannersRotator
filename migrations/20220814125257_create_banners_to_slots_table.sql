-- +goose Up
-- +goose StatementBegin
create table if not exists banners_to_slots
(
    banner_id int not null,
    slot_id int not null
);

create unique index banners_to_slots_banner_id_slot_id_uindex
    on banners_to_slots (banner_id, slot_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banners_to_slots;
-- +goose StatementEnd
