-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

create schema cmc;

create table cmc.currency
(
    id                                  bigint                  not null,
    symbol                              text                    not null,
    slug                                text                    not null,
    name                                text                    not null,
    is_for_observing                    bool                    not null,
    CONSTRAINT currency__id__pk PRIMARY KEY (id) include (symbol, slug, name, uri)
);


create table cmc.price_and_cap
(
    currency_id                         bigint                  not null,
    price                               double precision        not null,
    daily_volume                        double precision        not null,
    cap                                 double precision        not null,
    ts                                  timestamp               not null,
    CONSTRAINT price_and_cap__id__fk FOREIGN KEY (id) REFERENCES cmc.crypto(id)
);

create unique index price_and_cap__id__ts__pk ON cmc.price_and_cap (id, ts);

select public.create_hypertable('cmc.price_and_cap', 'ts', chunk_time_interval => INTERVAL '1 year');


create table cmc.concentration
(
    currency_id                         bigint                  not null,
    whales                              double precision        not null,
    investors                           double precision        not null,
    retail                              double precision        not null,
    others                              double precision        not null,
    d                                   date                    not null,
    CONSTRAINT concentration__id__fk FOREIGN KEY (id) REFERENCES cmc.crypto(id)
);




-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop schema cmc cascade;
-- +goose StatementEnd
