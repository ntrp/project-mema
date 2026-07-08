-- name: UpsertMediaFileFact :one
insert into app.media_file_facts (
    id,
    media_item_id,
    season_id,
    episode_id,
    file_path,
    quality_id,
    container_format,
    container_format_name,
    container_bitrate,
    duration_ms,
    size_bytes,
    source_kind,
    probed_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    sqlc.arg(probed_at)
)
on conflict (media_item_id, file_path) do update set
    season_id = excluded.season_id,
    episode_id = excluded.episode_id,
    quality_id = excluded.quality_id,
    container_format = excluded.container_format,
    container_format_name = excluded.container_format_name,
    container_bitrate = excluded.container_bitrate,
    duration_ms = excluded.duration_ms,
    size_bytes = excluded.size_bytes,
    source_kind = excluded.source_kind,
    probed_at = excluded.probed_at,
    updated_at = now()
returning id, media_item_id, season_id, episode_id, file_path, quality_id, container_format,
    container_format_name, container_bitrate, duration_ms, size_bytes, source_kind,
    probed_at, created_at, updated_at;

-- name: DeleteMediaFileTracksForFact :exec
delete from app.media_file_tracks
where media_file_fact_id = $1;

-- name: InsertMediaFileTrack :one
insert into app.media_file_tracks (
    id,
    media_file_fact_id,
    media_item_id,
    file_path,
    stream_index,
    track_type,
    language_id,
    codec,
    channels,
    bitrate_kbps,
    width,
    height,
    hdr_format,
    pixel_format,
    bit_depth,
    format,
    title,
    disposition
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    sqlc.arg(disposition)::jsonb
)
returning id, media_file_fact_id, media_item_id, file_path, stream_index, track_type,
    language_id, codec, channels, bitrate_kbps, width, height, hdr_format, pixel_format,
    bit_depth, format, title, disposition, created_at, updated_at;

-- name: ListMediaFileFactsForItem :many
select id, media_item_id, season_id, episode_id, file_path, quality_id, container_format,
    container_format_name, container_bitrate, duration_ms, size_bytes, source_kind,
    probed_at, created_at, updated_at
from app.media_file_facts
where media_item_id = $1
order by file_path;

-- name: ListMediaFileTracksForItem :many
select id, media_file_fact_id, media_item_id, file_path, stream_index, track_type,
    language_id, codec, channels, bitrate_kbps, width, height, hdr_format, pixel_format,
    bit_depth, format, title, disposition, created_at, updated_at
from app.media_file_tracks
where media_item_id = $1
order by file_path, stream_index, track_type;
