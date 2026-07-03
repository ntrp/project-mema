# Behavior Test Catalog

Behavior specs in this directory are the shared scenario pool for acceptance
coverage. Scenario IDs are stable traceability anchors:

- `SCN-AUTH-001`: domain prefix plus three-digit sequence.
- Tags declare the intended runner level: `@unit`, `@api`, `@integration`,
  `@e2e`.
- Tests assert observable inputs, outputs, and side effects. They should not
  depend on private implementation details.

Use Gherkin here for user and API behavior. Keep direct Go and TypeScript unit
tests close to the code, and reference the scenario ID in the test name or
comment when they cover a cataloged behavior.
