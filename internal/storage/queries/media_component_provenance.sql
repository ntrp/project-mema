-- name: UpsertMediaComponentProvenance :one
insert into app.media_component_provenance (
    id,
    media_item_id,
    component_type,
    component_key,
    release_group,
    release_name,
    release_id,
    source_provider,
    source_file_path,
    retained_source_id,
    source_stream_id,
    transformation_chain
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
    sqlc.arg(transformation_chain)::jsonb
)
on conflict (media_item_id, component_type, component_key) do update set
    release_group = excluded.release_group,
    release_name = excluded.release_name,
    release_id = excluded.release_id,
    source_provider = excluded.source_provider,
    source_file_path = excluded.source_file_path,
    retained_source_id = excluded.retained_source_id,
    source_stream_id = excluded.source_stream_id,
    transformation_chain = excluded.transformation_chain,
    updated_at = now()
returning id, media_item_id, component_type, component_key, release_group, release_name, release_id,
    source_provider, source_file_path, retained_source_id, source_stream_id, transformation_chain,
    created_at, updated_at;

-- name: ListMediaComponentProvenance :many
select id, media_item_id, component_type, component_key, release_group, release_name, release_id,
    source_provider, source_file_path, retained_source_id, source_stream_id, transformation_chain,
    created_at, updated_at
from app.media_component_provenance
where media_item_id = $1
order by component_type, component_key;
