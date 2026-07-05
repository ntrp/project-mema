create table if not exists river_job (
    id bigint primary key,
    state text not null,
    kind text not null,
    queue text not null,
    attempt integer not null,
    max_attempts integer not null,
    priority integer not null,
    args jsonb not null default '{}'::jsonb,
    metadata jsonb not null default '{}'::jsonb,
    errors jsonb[] not null default '{}'::jsonb[],
    scheduled_at timestamptz not null,
    created_at timestamptz not null,
    attempted_at timestamptz,
    finalized_at timestamptz
);
