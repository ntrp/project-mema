-- +goose Up
create table if not exists app.languages (
    code text primary key,
    display_name text not null,
    aliases jsonb not null default '[]'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint languages_code_check check (code = upper(code) and length(code) between 2 and 8),
    constraint languages_display_name_check check (length(trim(display_name)) > 0),
    constraint languages_aliases_array_check check (jsonb_typeof(aliases) = 'array')
);

create unique index if not exists idx_languages_display_name_lower
    on app.languages (lower(display_name));

-- +goose Down
drop table if exists app.languages;
