-- name: ListMediaComponentCompatibilityForSource :many
select id, media_item_id, base_source_id, component_source_id, confidence_state, automation_state, review_state, reason, runtime_delta_ms, evidence, review_reason, reviewed_at, created_at, updated_at
from app.media_component_compatibility_decisions
where component_source_id = sqlc.arg(component_source_id)
order by created_at desc, id;

-- name: GetMediaComponentCompatibility :one
select id, media_item_id, base_source_id, component_source_id, confidence_state, automation_state, review_state, reason, runtime_delta_ms, evidence, review_reason, reviewed_at, created_at, updated_at
from app.media_component_compatibility_decisions
where id = sqlc.arg(id)
  and media_item_id = sqlc.arg(media_item_id);

-- name: UpsertMediaComponentCompatibility :one
insert into app.media_component_compatibility_decisions (
    id,
    media_item_id,
    base_source_id,
    component_source_id,
    confidence_state,
    automation_state,
    review_state,
    reason,
    runtime_delta_ms,
    evidence
) values (
    sqlc.arg(id),
    sqlc.arg(media_item_id),
    sqlc.arg(base_source_id),
    sqlc.arg(component_source_id),
    sqlc.arg(confidence_state),
    sqlc.arg(automation_state),
    sqlc.arg(review_state),
    sqlc.arg(reason),
    sqlc.arg(runtime_delta_ms),
    sqlc.arg(evidence)
)
on conflict (media_item_id, base_source_id, component_source_id) do update
set confidence_state = excluded.confidence_state,
    automation_state = excluded.automation_state,
    review_state = case
        when app.media_component_compatibility_decisions.review_state in ('approved', 'rejected')
            then app.media_component_compatibility_decisions.review_state
        else excluded.review_state
    end,
    reason = excluded.reason,
    runtime_delta_ms = excluded.runtime_delta_ms,
    evidence = excluded.evidence,
    updated_at = now()
returning id, media_item_id, base_source_id, component_source_id, confidence_state, automation_state, review_state, reason, runtime_delta_ms, evidence, review_reason, reviewed_at, created_at, updated_at;

-- name: ReviewMediaComponentCompatibility :one
update app.media_component_compatibility_decisions
set review_state = sqlc.arg(review_state),
    automation_state = sqlc.arg(automation_state),
    review_reason = sqlc.arg(review_reason),
    reviewed_at = now(),
    updated_at = now()
where id = sqlc.arg(id)
  and media_item_id = sqlc.arg(media_item_id)
  and component_source_id = sqlc.arg(component_source_id)
returning id, media_item_id, base_source_id, component_source_id, confidence_state, automation_state, review_state, reason, runtime_delta_ms, evidence, review_reason, reviewed_at, created_at, updated_at;
