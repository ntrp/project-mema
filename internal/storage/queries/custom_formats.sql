-- name: ListCustomFormats :many
select id, name, include_in_rename_template, include_specs, exclude_specs, created_at, updated_at
from app.custom_formats
order by lower(name);

-- name: CreateCustomFormat :one
insert into app.custom_formats (id, name, include_in_rename_template, include_specs, exclude_specs)
values ($1, $2, $3, sqlc.arg(include_specs)::jsonb, sqlc.arg(exclude_specs)::jsonb)
returning id, name, include_in_rename_template, include_specs, exclude_specs, created_at, updated_at;

-- name: UpdateCustomFormat :one
update app.custom_formats
set name = $2,
    include_in_rename_template = $3,
    include_specs = sqlc.arg(include_specs)::jsonb,
    exclude_specs = sqlc.arg(exclude_specs)::jsonb,
    updated_at = now()
where id = $1
returning id, name, include_in_rename_template, include_specs, exclude_specs, created_at, updated_at;

-- name: DeleteCustomFormat :execrows
delete from app.custom_formats
where id = $1;
