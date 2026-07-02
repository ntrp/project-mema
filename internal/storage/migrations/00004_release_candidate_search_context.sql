-- +goose Up
alter table app.media_release_candidates
	add column if not exists search_kind text not null default 'title',
	add column if not exists requested_season integer,
	add column if not exists requested_episode integer;
