create extension if not exists pgcrypto;

-- таблица с пользователями
create table if not exists users (
    id bigserial primary key,
    login text not null unique,
    password_hash text not null,
    created_at timestamptz not null default now()
);

-- таблица сессий
create table if not exists sessions (
    id text primary key default encode(gen_random_bytes(16),'hex'),
    uid bigint not null, -- user id
    created_at timestamptz not null default now()
);

-- таблица с файлами
create table if not exists assets (
    name text not null,
    uid bigint not null, -- user id
    data bytea not null,
    created_at timestamptz not null default now(),
    primary key (name, uid)
);

-- тестовый пользователь
insert into users
    (login, password_hash)
values
    ('alice', encode(digest('secret', 'md5'),'hex'))
on conflict do nothing;