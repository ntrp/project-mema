-- name: ListLanguages :many
select code, display_name, aliases, created_at, updated_at
from app.languages
order by lower(display_name), code;

-- name: UpsertLanguage :one
insert into app.languages (code, display_name, aliases)
values ($1, $2, sqlc.arg(aliases)::jsonb)
on conflict (code) do update
set display_name = excluded.display_name,
    aliases = excluded.aliases,
    updated_at = now()
returning code, display_name, aliases, created_at, updated_at;

-- name: DeleteLanguage :execrows
delete from app.languages
where code = $1;
