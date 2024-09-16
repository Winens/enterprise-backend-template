create unlogged table if not exists sessions (
    id uuid primary key,
    user_id bigint not null,
    foreign key (user_id) references users(id) on delete cascade,

    ip varchar(45),
    user_agent varchar(255),

    created_at timestamptz default current_timestamp,
    last_seen_at timestamptz default current_timestamp
);
