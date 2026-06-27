create schema if not exists app;

create table if not exists app.users (
    id uuid primary key,
    username text not null unique,
    password_hash text not null,
    role text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists app.sessions (
    id text primary key,
    user_id uuid not null references app.users(id) on delete cascade,
    expires_at timestamptz not null,
    created_at timestamptz not null default now()
);

create table if not exists app.download_clients (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('transmission', 'sabnzbd')),
    base_url text not null,
    username text,
    password text,
    api_key text,
    category text,
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_download_clients_priority
    on app.download_clients (priority, name);

create table if not exists app.indexers (
    id uuid primary key,
    name text not null,
    type text not null check (type in ('torznab', 'newznab', 'rss')),
    base_url text not null,
    api_key text,
    categories integer[] not null default '{}',
    enabled boolean not null default true,
    priority integer not null default 100 check (priority >= 0 and priority <= 1000),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index if not exists idx_indexers_priority
    on app.indexers (priority, name);
