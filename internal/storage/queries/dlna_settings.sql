-- name: GetDLNASettings :one
insert into app.dlna_settings (
    id,
    enabled,
    friendly_name,
    interfaces,
    allowed_cidrs,
    announce_interval_seconds,
    transcode_enabled,
    thumbnails_enabled,
    subtitles_enabled,
    default_renderer_profile,
    device_uuid
)
values (true, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
on conflict (id) do update set id = excluded.id
returning
    enabled,
    friendly_name,
    interfaces,
    allowed_cidrs,
    announce_interval_seconds,
    transcode_enabled,
    thumbnails_enabled,
    subtitles_enabled,
    default_renderer_profile,
    device_uuid,
    created_at,
    updated_at;

-- name: UpdateDLNASettings :one
insert into app.dlna_settings (
    id,
    enabled,
    friendly_name,
    interfaces,
    allowed_cidrs,
    announce_interval_seconds,
    transcode_enabled,
    thumbnails_enabled,
    subtitles_enabled,
    default_renderer_profile,
    device_uuid
)
values (true, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
on conflict (id) do update
set enabled = excluded.enabled,
    friendly_name = excluded.friendly_name,
    interfaces = excluded.interfaces,
    allowed_cidrs = excluded.allowed_cidrs,
    announce_interval_seconds = excluded.announce_interval_seconds,
    transcode_enabled = excluded.transcode_enabled,
    thumbnails_enabled = excluded.thumbnails_enabled,
    subtitles_enabled = excluded.subtitles_enabled,
    default_renderer_profile = excluded.default_renderer_profile,
    updated_at = now()
returning
    enabled,
    friendly_name,
    interfaces,
    allowed_cidrs,
    announce_interval_seconds,
    transcode_enabled,
    thumbnails_enabled,
    subtitles_enabled,
    default_renderer_profile,
    device_uuid,
    created_at,
    updated_at;
