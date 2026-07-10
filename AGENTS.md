# Agents instructions

## Svelte AI tools

This project uses Svelte 5 (`web/package.json`: `svelte` ^5.56.1). Do not target Svelte 4 or legacy Svelte syntax for new code.

You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Use these tools for Svelte work:

### Available Svelte MCP tools

1. `list-sections` / `svelte_list-sections`
   - Use this first to discover available documentation sections. It returns titles, use cases, and paths.
   - When asked about Svelte or SvelteKit topics, use it at the start of the chat to find relevant sections.

2. `get-documentation` / `svelte_get-documentation`
   - Retrieves full documentation for specific sections.
   - After `list-sections`, analyze the returned sections, especially `use_cases`, then fetch all docs relevant to the user's task.

3. `svelte-autofixer` / `svelte_svelte-autofixer`
   - Analyzes Svelte code and returns issues and suggestions.
   - MUST be used whenever writing or editing Svelte components or Svelte modules before returning code. Always target Svelte 5. Keep fixing and rerunning until no issues or suggestions remain.

4. `playground-link` / `svelte_playground-link`
   - Generates a Svelte Playground link with provided code.
   - After completing code that was not written to project files, ask the user if they want a playground link. Only call this tool after user confirmation.

For `.svelte`, `.svelte.ts`, and `.svelte.js` files, prefer the project `svelte-file-editor` / `svelte-code-writer` skills when available. If MCP tools are unavailable, use the CLI fallback:

```bash
npx -y @sveltejs/mcp list-sections
npx -y @sveltejs/mcp get-documentation "<section1>,<section2>"
npx -y @sveltejs/mcp svelte-autofixer "<code_or_path>" --svelte-version 5
```

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

- if the database schema or seed changes, reset the db, dev seed, local seed
