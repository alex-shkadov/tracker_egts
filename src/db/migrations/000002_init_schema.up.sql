CREATE TABLE IF NOT EXISTS trackers (
    id serial2 not null primary key,
    title varchar(255),
    imei varchar(15),
    transport_number varchar(255),
    description TEXT,
    is_active bool not null default false,
    created_at timestamp with time zone not null default now()::timestamp,
    updated_at timestamp with time zone default null,
    deleted_at timestamp with time zone default null,
    constraint trackers_imei_uniq UNIQUE (imei)
);

CREATE TABLE IF NOT EXISTS service_data_records (
    id serial8 not null primary key,
    packet_id int2 NOT NULL,
    tracker_id int2 NOT NULL REFERENCES trackers(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    record_number int4 not null,
    object_identifier int4 not null,
    created_at timestamp with time zone not null default now()::timestamp,
    updated_at timestamp with time zone default null,
    deleted_at timestamp with time zone default null
);

CREATE TABLE IF NOT EXISTS sr_pos_data (
    id serial8 not null primary key,
    service_data_record_id int4 NOT NULL REFERENCES service_data_records (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    ntm timestamptz not null,
    latitude float not null,
    longitude float not null,
    mv bool default false,
    bb bool default false,
    spd int2 not null,
    alts int4 not null,
    dir int not null,
    dirh int not null,
    odm int4 not null,
    satellites int not null,
    record_number int2 not null,
    created_at timestamp with time zone not null default now()::timestamp,
    updated_at timestamp with time zone default null,
    deleted_at timestamp with time zone default null
);
