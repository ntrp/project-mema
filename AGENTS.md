# Agents instructions

- Always keep modules small and refactor them when they grow too much, this includes svelte components and go modules
- no backend module should be over 300 lines of code, excluding generated files
- no frontend component should be over 200 lines of code
- DRY, don't repeat yourself, if you find yourself copy pasting code, refactor it into a module or component
- Single responsibility principle, each module or component should have a single responsibility or a very small number of them

## Dev workflow

- When working on tasks start the frontend and backend
- Always bind dev servers to `0.0.0.0`; use `ADDR=0.0.0.0:18080 make dev-api` for the backend and `pnpm exec vite dev --host 0.0.0.0 --port 15173` from `web/` for the frontend
- if the backend changes restart it
- if the database schema or seed changes, reset the db and restart the backend
