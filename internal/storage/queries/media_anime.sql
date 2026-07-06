-- name: ListMediaProviderMappings :many
select *
from app.media_provider_mappings
where media_item_id = $1
order by canonical desc, provider_name asc, entity_type asc, external_id asc;

-- name: UpsertMediaProviderMapping :one
insert into app.media_provider_mappings (
    id,
    media_item_id,
    season_id,
    episode_id,
    entity_type,
    provider_name,
    provider_entity_type,
    external_id,
    canonical,
    confidence,
    source
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.narg(episode_id),
    sqlc.arg(entity_type),
    sqlc.arg(provider_name),
    sqlc.arg(provider_entity_type),
    sqlc.arg(external_id),
    sqlc.arg(canonical),
    sqlc.narg(confidence),
    sqlc.arg(source)::jsonb
)
on conflict (
    media_item_id,
    (coalesce(season_id, '00000000-0000-0000-0000-000000000000'::uuid)),
    (coalesce(episode_id, '00000000-0000-0000-0000-000000000000'::uuid)),
    provider_name,
    provider_entity_type,
    external_id
) do update
set canonical = excluded.canonical,
    confidence = excluded.confidence,
    source = excluded.source,
    updated_at = now()
returning *;

-- name: ListMediaItemAliases :many
select *
from app.media_item_aliases
where media_item_id = $1
order by alias_kind asc, alias asc;

-- name: UpsertMediaItemAlias :one
insert into app.media_item_aliases (
    id,
    media_item_id,
    alias,
    normalized_alias,
    language,
    alias_kind,
    provider_name,
    provider_mapping_id,
    source
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(alias),
    sqlc.arg(normalized_alias),
    sqlc.narg(language),
    sqlc.arg(alias_kind),
    sqlc.narg(provider_name),
    sqlc.narg(provider_mapping_id),
    sqlc.arg(source)::jsonb
)
on conflict (
    media_item_id,
    normalized_alias,
    alias_kind,
    (coalesce(language, '')),
    (coalesce(provider_name, ''))
) do update
set alias = excluded.alias,
    provider_mapping_id = excluded.provider_mapping_id,
    source = excluded.source,
    updated_at = now()
returning *;

-- name: ListMediaEpisodeNumbering :many
select *
from app.media_episode_numbering
where media_item_id = $1
order by provider_name asc, numbering_scheme asc, season_number asc, episode_number asc, absolute_number asc;

-- name: UpsertMediaEpisodeNumbering :one
insert into app.media_episode_numbering (
    id,
    media_item_id,
    season_id,
    episode_id,
    provider_name,
    numbering_scheme,
    season_number,
    episode_number,
    absolute_number,
    source
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.narg(season_id),
    sqlc.arg(episode_id),
    sqlc.arg(provider_name),
    sqlc.arg(numbering_scheme),
    sqlc.narg(season_number),
    sqlc.narg(episode_number),
    sqlc.narg(absolute_number),
    sqlc.arg(source)::jsonb
)
on conflict (episode_id, provider_name, numbering_scheme) do update
set season_id = excluded.season_id,
    season_number = excluded.season_number,
    episode_number = excluded.episode_number,
    absolute_number = excluded.absolute_number,
    source = excluded.source,
    updated_at = now()
returning *;
