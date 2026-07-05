# Agents instructions

- Always keep modules small and refactor them when they grow too much, this includes svelte components and go modules
- no backend module should be over 300 lines of code, excluding generated files
- no frontend component or module should be over 200 lines of code
- when multiple components are in the same feature organize them in a directory, try to keep maximum number files in a single directory to 10
- DRY, don't repeat yourself, if you find yourself copy pasting code, refactor it into a module or component
- Single responsibility principle, each module or component should have a single responsibility or a very small number of them
- Always use a modal for confirmation dialogs, never use a browser confirm dialog
- Always use a tooltip component for showing tooltips, never use a browser tooltip
- This project is not released yet; do not create new database migrations. Apply database schema changes directly to `internal/storage/migrations/00001_initial_schema.sql`, then reset the development database.
- never touch the dev.local.sql seed unless explicitly asked to update it

## Dev workflow

- When working on tasks start the frontend and backend
- Always bind dev servers to `0.0.0.0`;
  - use `ADDR=0.0.0.0:18080 make dev-api` for the backend
  - `NVIM_LISTEN_ADDRESS=/tmp/project-mema.nvim LAUNCH_EDITOR=/Users/ntrp/_pws/project-mema/scripts/open-in-nvim.sh pnpm exec vite dev --host 0.0.0.0 --port 15173` in /web for the frontend
- if the backend changes restart it
- if the database schema or seed changes, reset the db and restart the backend
