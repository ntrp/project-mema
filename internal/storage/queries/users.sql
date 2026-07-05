-- name: EnsureDefaultAdminUser :exec
insert into app.users (id, username, password_hash, role)
select $1, $2, $3, 'admin'
where not exists (select 1 from app.users where username = $2);

-- name: ListUsers :many
select id, username, password_hash, display_name, picture_url, role, created_at, updated_at
from app.users
order by username asc;

-- name: GetUser :one
select id, username, password_hash, display_name, picture_url, role, created_at, updated_at
from app.users
where id = $1;

-- name: GetUserByUsername :one
select id, username, password_hash, display_name, picture_url, role, created_at, updated_at
from app.users
where lower(username) = lower($1);

-- name: CreateUser :one
insert into app.users (id, username, password_hash, role)
values ($1, $2, $3, $4)
returning id, username, password_hash, display_name, picture_url, role, created_at, updated_at;

-- name: UpdateUser :one
update app.users
set username = $2,
    role = $3,
    updated_at = now()
where id = $1
returning id, username, password_hash, display_name, picture_url, role, created_at, updated_at;

-- name: UpdateUserWithPassword :one
update app.users
set username = $2,
    password_hash = $3,
    role = $4,
    updated_at = now()
where id = $1
returning id, username, password_hash, display_name, picture_url, role, created_at, updated_at;

-- name: DeleteUser :execrows
delete from app.users
where id = $1;

-- name: CountAdminUsers :one
select count(*)::int
from app.users
where role = 'admin';
