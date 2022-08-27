-- +goose Up
-- +goose StatementBegin
insert into banners (id, description)
values (1, 'azino 777'),
       (2, '1xbet'),
       (3, 'online job without experience'),
       (4, 'another scum');

insert into slots (id, description)
values (1, 'top'),
       (2, 'left'),
       (3, 'right'),
       (4, 'bottom');

insert into dem_groups (id, description)
values (1, 'male 18-25'),
       (2, 'male 25-40'),
       (3, 'male 41+'),
       (4, 'female 18-25'),
       (5, 'female 25-40'),
       (6, 'female 41+');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table banners;
truncate table slots;
truncate table dem_groups;
-- +goose StatementEnd
