---
name: svelte-code-writer
description: CLI tools for Svelte 5 documentation lookup and code analysis. MUST be used whenever creating, editing or analyzing any Svelte component (.svelte) or Svelte module (.svelte.ts/.svelte.js). If possible, this skill should be executed within the svelte-file-editor agent for optimal results.
---

# Svelte 5 Code Writer

This project uses Svelte 5. Always target Svelte 5 and prefer runes-mode patterns for new code. Do not use Svelte 4 examples or legacy syntax unless the user explicitly asks for migration/debugging of existing legacy code.

## CLI Tools

You have access to `@sveltejs/mcp` CLI for Svelte-specific assistance. Use these commands via `npx`:

### List Documentation Sections

```bash
npx -y @sveltejs/mcp list-sections
```

Lists all available Svelte 5 and SvelteKit documentation sections with titles and paths.

### Get Documentation

```bash
npx -y @sveltejs/mcp get-documentation "<section1>,<section2>,..."
```

Retrieves full documentation for specified sections. Use after `list-sections` to fetch relevant docs.

**Example:**

```bash
npx -y @sveltejs/mcp get-documentation "$state,$derived,$effect"
```

### Svelte Autofixer

```bash
npx -y @sveltejs/mcp svelte-autofixer "<code_or_path>" --svelte-version 5
```

Analyzes Svelte code and suggests fixes for common issues.

**Options:**

- `--async` - Enable async Svelte mode (default: false)
- `--svelte-version` - Target version. Use `5` for this project.

**Examples:**

```bash
# Analyze inline code (escape $ as \$)
npx -y @sveltejs/mcp svelte-autofixer '<script>let count = \$state(0);</script>'

# Analyze a file, targeting Svelte 5
npx -y @sveltejs/mcp svelte-autofixer ./src/lib/Component.svelte --svelte-version 5
```

**Important:** When passing code with runes (`$state`, `$derived`, etc.) via the terminal, escape the `$` character as `\$` to prevent shell variable substitution.

## Workflow

1. **Uncertain about syntax?** Run `list-sections` then `get-documentation` for relevant topics
2. **Reviewing/debugging?** Run `svelte-autofixer` on the code to detect issues
3. **Always validate** - Run `svelte-autofixer --svelte-version 5` before finalizing any Svelte component
