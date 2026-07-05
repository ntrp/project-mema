-- name: CreateSession :exec
insert into app.sessions (id, user_id, expires_at)
values ($1, $2, $3)
on conflict (id) do update
set user_id = excluded.user_id,
    expires_at = excluded.expires_at,
    updated_at = now();

-- name: GetSession :one
select s.id,
    s.user_id,
    s.expires_at,
    s.created_at,
    s.updated_at,
    u.username,
    u.display_name,
    u.picture_url,
    u.role
from app.sessions s
join app.users u on u.id = s.user_id
where s.id = $1
    and s.expires_at > $2;

-- name: DeleteSession :execrows
delete from app.sessions
where id = $1;

-- name: DeleteExpiredSessions :execrows
delete from app.sessions
where expires_at <= $1;
