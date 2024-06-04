-- +goose Up
create table events (
    id serial primary key,
    owner_id bigint not null,
    title text not null,
    description text,
    begin_datetime timestamp with time zone not null,
    duration bigint,
    notify bigint
);

create index owner_idx on events (owner_id);
create index start_idx on events using btree (begin_datetime);

-- +goose Down
drop table events;