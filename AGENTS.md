# Agents instructions

- Always keep modules small and refactor them when they grow too much, this includes svelte components and go modules
- no backend module should be over 300 lines of code, excluding generated files
- no frontend component should be over 200 lines of code

## Dev workflow

- When working on tasks start the frontend and backend
- if the backend changes restart it
- if the database schema or seed changes, reset the db and restart the backend
