-- name: ListTags :many
select id, name, created_at, updated_at
from app.tags
order by lower(name);

-- name: UpsertTagByName :one
insert into app.tags (id, name)
values ($1, $2)
on conflict (lower(name)) do update
set name = excluded.name, updated_at = now()
returning id, name, created_at, updated_at;

-- name: UpdateTag :one
update app.tags
set name = $2, updated_at = now()
where id = $1
returning id, name, created_at, updated_at;

-- name: DeleteTag :execrows
delete from app.tags
where id = $1;

-- name: EnsureTag :one
insert into app.tags (id, name)
values ($1, $2)
on conflict (lower(name)) do update
set name = excluded.name, updated_at = now()
returning id;

-- name: DeleteMediaItemTags :exec
delete from app.media_item_tags
where media_item_id = $1;

-- name: AddMediaItemTag :exec
insert into app.media_item_tags (media_item_id, tag_id)
values ($1, $2)
on conflict do nothing;

-- name: DeleteMediaRequestTags :exec
delete from app.media_request_tags
where media_request_id = $1;

-- name: AddMediaRequestTag :exec
insert into app.media_request_tags (media_request_id, tag_id)
values ($1, $2)
on conflict do nothing;
