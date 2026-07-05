-- name: GetDatabaseVersion :one
select current_setting('server_version')::text;
