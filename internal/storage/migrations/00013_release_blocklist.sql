-- +goose Up
create table if not exists app.release_blocklist (
    id uuid primary key,
    media_item_id uuid not null references app.media_items(id) on delete cascade,
    release_title text not null,
    indexer_name text not null,
    indexer_type text not null default '',
    download_url text,
    info_url text,
    guid text,
    reason text not null,
    source text not null,
    temporary boolean not null default true,
    expires_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint release_blocklist_expiry_check check (
        (temporary = false and expires_at is null) or (temporary = true and expires_at is not null)
    )
);

create index if not exists idx_release_blocklist_media
    on app.release_blocklist (media_item_id, created_at desc);

create index if not exists idx_release_blocklist_expiry
    on app.release_blocklist (temporary, expires_at);

-- +goose Down
drop table if exists app.release_blocklist;
