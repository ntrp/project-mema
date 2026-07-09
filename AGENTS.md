# Agents instructions

- no backend module should be over 300 lines of code, excluding generated files
- no frontend component or module should be over 200 lines of code excluding generated files
- when multiple components are in the same feature organize them in a directory, keep maximum number files in a single directory to 10
- DRY, don't repeat yourself, if you find yourself copy pasting code, refactor it into a module or component
- Single responsibility principle, each module or component should have a single responsibility or a very small number of them
- Always use a modal for confirmation dialogs, never use a browser confirm dialog
- Always use a tooltip component for showing tooltips, never use a browser tooltip
- This project is not released yet; do not create new database migrations. Apply database schema changes directly to `internal/storage/migrations/00001_initial_schema.sql`, then reset the development database.
- never touch the dev.local.sql seed unless explicitly asked to update it
- keep coverage more than 60% on every file, only data objects and non business logic code can be excluded
- for every new feature keep the docs sections up to date: update user guide for user facing changes, update architecture for internal changes

## Dev workflow

- When working on tasks start the backend
- Start backend with ADDR=0.0.0.0:18080 make dev-api` for the backend
- if the backend changes restart it
- if the database schema or seed changes, reset the db, dev seed, local seed and restart the backend
