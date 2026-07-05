-- name: BlockRelease :one
insert into app.release_blocklist (
    id,
    media_item_id,
    release_title,
    indexer_name,
    indexer_protocol,
    download_client_name,
    download_url,
    info_url,
    guid,
    reason,
    source,
    temporary,
    expires_at
)
values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(release_title),
    sqlc.arg(indexer_name),
    sqlc.arg(indexer_protocol),
    sqlc.arg(download_client_name),
    sqlc.narg(download_url),
    sqlc.narg(info_url),
    sqlc.narg(guid),
    sqlc.arg(reason),
    sqlc.arg(source),
    sqlc.arg(temporary),
    sqlc.narg(expires_at)
)
on conflict (id) do update set updated_at = now()
returning id,
    media_item_id,
    ''::text as media_title,
    ''::text as media_type,
    release_title,
    indexer_name,
    coalesce(nullif(indexer_protocol, ''), 'torrent')::text as indexer_protocol,
    download_client_name,
    download_url,
    info_url,
    guid,
    reason,
    source,
    temporary,
    expires_at,
    created_at,
    updated_at;

-- name: ListReleaseBlocklist :many
select b.id,
    b.media_item_id,
    m.title as media_title,
    m.media_type,
    b.release_title,
    b.indexer_name,
    coalesce(nullif(b.indexer_protocol, ''), i.protocol, 'torrent')::text as indexer_protocol,
    b.download_client_name,
    b.download_url,
    b.info_url,
    b.guid,
    b.reason,
    b.source,
    b.temporary,
    b.expires_at,
    b.created_at,
    b.updated_at
from app.release_blocklist b
join app.media_items m on m.id = b.media_item_id
left join app.indexers i on lower(i.name) = lower(b.indexer_name)
order by b.created_at desc
limit 200;

-- name: CleanupExpiredReleaseBlocks :execrows
delete from app.release_blocklist
where temporary = true and expires_at <= now();

-- name: DeleteReleaseBlocklistItem :execrows
delete from app.release_blocklist
where id = $1;

-- name: ClearReleaseBlocklist :execrows
delete from app.release_blocklist;

-- name: FindReleaseBlock :one
select b.id,
    b.media_item_id,
    ''::text as media_title,
    ''::text as media_type,
    b.release_title,
    b.indexer_name,
    coalesce(nullif(b.indexer_protocol, ''), 'torrent')::text as indexer_protocol,
    b.download_client_name,
    b.download_url,
    b.info_url,
    b.guid,
    b.reason,
    b.source,
    b.temporary,
    b.expires_at,
    b.created_at,
    b.updated_at
from app.release_blocklist b
where b.media_item_id = sqlc.arg(media_item_id)
    and (b.temporary = false or b.expires_at > now())
    and (
        (sqlc.narg(guid)::text is not null and b.guid = sqlc.narg(guid)::text)
        or (sqlc.narg(info_url)::text is not null and b.info_url = sqlc.narg(info_url)::text)
        or (sqlc.narg(download_url)::text is not null and b.download_url = sqlc.narg(download_url)::text)
        or lower(b.release_title) = lower(sqlc.arg(title))
    )
order by b.created_at desc
limit 1;
