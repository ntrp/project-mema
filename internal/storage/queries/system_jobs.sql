-- name: ListSystemJobs :many
select id,
    state::text as state,
    kind,
    queue,
    attempt::int as attempt,
    max_attempts::int as max_attempts,
    priority::int as priority,
    args::text as args,
    metadata::text as metadata,
    coalesce(array_to_json(errors), '[]'::json)::text as errors,
    coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text)::text as info_message,
    scheduled_at,
    created_at,
    attempted_at,
    finalized_at
from river_job
where (cardinality(sqlc.arg(states)::text[]) = 0 or state::text = any(sqlc.arg(states)::text[]))
    and (sqlc.arg(queue)::text = '' or queue = sqlc.arg(queue)::text)
    and (sqlc.arg(kind)::text = '' or kind = sqlc.arg(kind)::text)
    and (
        sqlc.arg(search_query)::text = ''
        or kind ilike '%' || sqlc.arg(search_query)::text || '%'
        or queue ilike '%' || sqlc.arg(search_query)::text || '%'
        or args::text ilike '%' || sqlc.arg(search_query)::text || '%'
        or errors::text ilike '%' || sqlc.arg(search_query)::text || '%'
    )
order by coalesce(finalized_at, attempted_at, scheduled_at, created_at) desc, id desc
limit sqlc.arg(row_limit);

-- name: GetSystemJob :one
select id,
    state::text as state,
    kind,
    queue,
    attempt::int as attempt,
    max_attempts::int as max_attempts,
    priority::int as priority,
    args::text as args,
    metadata::text as metadata,
    coalesce(array_to_json(errors), '[]'::json)::text as errors,
    coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text)::text as info_message,
    scheduled_at,
    created_at,
    attempted_at,
    finalized_at
from river_job
where id = $1;
