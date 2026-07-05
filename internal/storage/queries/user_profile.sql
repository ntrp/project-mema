-- name: GetUserProfile :one
select id, username, display_name, picture_url, role, updated_at
from app.users
where id = $1;

-- name: UpdateUserProfile :one
update app.users
set display_name = $2,
    picture_url = $3,
    updated_at = now()
where id = $1
returning id, username, display_name, picture_url, role, updated_at;
