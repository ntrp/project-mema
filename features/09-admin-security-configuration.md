# PRD: Administration, Security, and Configuration

Status: Draft

## Summary

Mema is self-hosted software that needs reliable configuration, authentication options, secrets handling, observability, backups, and safe administrative controls.

## Goals

- Easy local deployment.
- Safe storage of credentials and API keys.
- Clear health checks.
- Backup and restore.
- Observable background jobs.
- Optional authentication for LAN-only users and stronger authentication for exposed deployments.

## Functional Requirements

- Users can configure root folders, clients, indexers, providers, profiles, and naming.
- Users can test each external integration.
- Users can view system health.
- Users can view queue, history, logs, and worker state.
- Users can backup and restore configuration and database.
- Users can configure authentication.
- Users can configure API keys.
- Users can configure webhooks and notifications.
- Users can configure scheduled tasks.
- Users can configure external tool paths.
- Users can configure retention and cleanup policies.

## Deployment Requirements

- Docker image.
- Docker Compose example.
- Single Go binary serving the SvelteKit build, if feasible.
- PostgreSQL support from the start.
- No migration tooling until the project is live.
- Configurable data directory.
- Configurable temp and work directories for muxing.

## Security Requirements

- Secrets must not be logged.
- API keys and provider credentials must be stored safely.
- Dangerous file operations require explicit roots.
- Path traversal must be prevented.
- External commands must be invoked with controlled arguments.
- Authentication should support at least local username/password or reverse-proxy auth.

## Acceptance Criteria

- A user can deploy Mema with Docker Compose.
- A user can configure at least one root folder and one download client.
- Integration tests show clear pass/fail status.
- Logs redact secrets.
- Backups include configuration and database.
- Restore can recover profiles and library state.

## Open Questions

- PostgreSQL is required from day one.
- Migration tooling is deferred until the project is live.
- Should auth be mandatory or optional?
- Should reverse proxy header auth be supported?
- Should OAuth/OIDC be supported?
- Should there be multi-user roles?
- Should secrets be encrypted at rest, and how should the key be managed?
- Should the app expose Prometheus metrics?
- Should logs be plain text, JSON, or configurable?
- Should background jobs survive restart with persisted queues?
- Should the UI support dark mode by default?
