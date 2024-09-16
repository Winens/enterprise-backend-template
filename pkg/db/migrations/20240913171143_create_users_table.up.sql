create table if not exists users (
    id bigserial primary key,

    first_name varchar(64) not null,
    last_name varchar(64) not null,

    email varchar(255) not null unique,
    email_confirmed boolean not null default false,
    email_confirmed_at timestamptz null,

    password_hash varchar(72) not null,

    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);
