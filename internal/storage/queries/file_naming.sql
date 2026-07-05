-- name: EnsureDefaultFileNamingSettings :exec
insert into app.file_naming_settings (
    id,
    movie_file_format,
    movie_folder_format,
    series_episode_format,
    daily_episode_format,
    anime_episode_format,
    series_folder_format,
    season_folder_format,
    specials_folder_format
)
values (1, $1, $2, $3, $4, $5, $6, $7, $8)
on conflict do nothing;

-- name: GetFileNamingSettings :one
select movie_file_format,
    movie_folder_format,
    series_episode_format,
    daily_episode_format,
    anime_episode_format,
    series_folder_format,
    season_folder_format,
    specials_folder_format,
    created_at,
    updated_at
from app.file_naming_settings
where id = 1;

-- name: UpdateFileNamingSettings :one
update app.file_naming_settings
set movie_file_format = $1,
    movie_folder_format = $2,
    series_episode_format = $3,
    daily_episode_format = $4,
    anime_episode_format = $5,
    series_folder_format = $6,
    season_folder_format = $7,
    specials_folder_format = $8,
    updated_at = now()
where id = 1
returning movie_file_format,
    movie_folder_format,
    series_episode_format,
    daily_episode_format,
    anime_episode_format,
    series_folder_format,
    season_folder_format,
    specials_folder_format,
    created_at,
    updated_at;
