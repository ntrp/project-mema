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
