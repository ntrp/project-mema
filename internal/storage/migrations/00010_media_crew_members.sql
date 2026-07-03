-- +goose Up
alter table app.media_items
    add column if not exists crew_members jsonb not null default '[]'::jsonb;

-- +goose StatementBegin
do $$
begin
    alter table app.media_items drop constraint if exists media_items_crew_members_check;
    alter table app.media_items
        add constraint media_items_crew_members_check check (jsonb_typeof(crew_members) = 'array');
end $$;
-- +goose StatementEnd

-- +goose Down
alter table app.media_items
    drop column if exists crew_members;
