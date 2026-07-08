-- name: ListDLNARendererProfiles :many
select *
from app.dlna_renderer_profiles
order by priority, name;

-- name: GetDLNARendererProfile :one
select *
from app.dlna_renderer_profiles
where id = $1;

-- name: CreateDLNARendererProfile :one
insert into app.dlna_renderer_profiles (
    id,
    name,
    vendor,
    device_class,
    source,
    source_version,
    customized,
    enabled,
    priority,
    icon_key,
    notes,
    match_rules,
    capability_rules,
    delivery_settings,
    dlna_flags,
    subtitle_rules,
    artwork_rules,
    metadata_rules,
    quirks
)
values (
    $1, $2, $3, $4, 'user', 1, true, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
returning *;

-- name: UpdateDLNARendererProfile :one
update app.dlna_renderer_profiles
set name = $2,
    vendor = $3,
    device_class = $4,
    enabled = $5,
    priority = $6,
    icon_key = $7,
    notes = $8,
    match_rules = $9,
    capability_rules = $10,
    delivery_settings = $11,
    dlna_flags = $12,
    subtitle_rules = $13,
    artwork_rules = $14,
    metadata_rules = $15,
    quirks = $16,
    source = 'user',
    customized = true,
    updated_at = now()
where id = $1
returning *;

-- name: DeleteDLNARendererProfile :execrows
delete from app.dlna_renderer_profiles
where id = $1;

-- name: ResetDLNARendererProfile :one
update app.dlna_renderer_profiles as profile
set name = defaults.name,
    vendor = defaults.vendor,
    device_class = defaults.device_class,
    source = 'mema_seed',
    source_version = defaults.source_version,
    customized = false,
    enabled = defaults.enabled,
    priority = defaults.priority,
    icon_key = defaults.icon_key,
    notes = defaults.notes,
    match_rules = defaults.match_rules,
    capability_rules = defaults.capability_rules,
    delivery_settings = defaults.delivery_settings,
    dlna_flags = defaults.dlna_flags,
    subtitle_rules = defaults.subtitle_rules,
    artwork_rules = defaults.artwork_rules,
    metadata_rules = defaults.metadata_rules,
    quirks = defaults.quirks,
    updated_at = now()
from app.dlna_renderer_profile_defaults as defaults
where profile.id = $1
  and defaults.id = profile.id
returning profile.*;

-- name: RebaseDLNARendererProfile :one
update app.dlna_renderer_profiles as profile
set source = 'user',
    source_version = defaults.source_version,
    customized = true,
    updated_at = now()
from app.dlna_renderer_profile_defaults as defaults
where profile.id = $1
  and defaults.id = profile.id
returning profile.*;

-- name: ListDLNARendererDeviceOverrides :many
select *
from app.dlna_renderer_device_overrides
order by display_name, renderer_uuid, ip_address;

-- name: UpsertDLNARendererDeviceOverride :one
insert into app.dlna_renderer_device_overrides (
    id,
    renderer_uuid,
    ip_address,
    profile_id,
    display_name,
    allowed,
    delivery_policy_overrides,
    notes
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (id) do update
set renderer_uuid = excluded.renderer_uuid,
    ip_address = excluded.ip_address,
    profile_id = excluded.profile_id,
    display_name = excluded.display_name,
    allowed = excluded.allowed,
    delivery_policy_overrides = excluded.delivery_policy_overrides,
    notes = excluded.notes,
    updated_at = now()
returning *;

-- name: DeleteDLNARendererDeviceOverride :execrows
delete from app.dlna_renderer_device_overrides
where id = $1;
