---
title: API Contract
description: OpenAPI ownership and generated client artifacts.
---

The API contract is source-controlled at:

```txt
api/openapi.yaml
```

The backend handler types are generated under:

```txt
internal/httpapi/openapi.gen.go
```

The frontend schema types are generated under:

```txt
web/src/lib/api/generated/schema.d.ts
```

## Contract Workflow

After editing `api/openapi.yaml`, regenerate both sides:

```sh
make api-generate
```

Before finishing contract work, verify generated artifacts:

```sh
make verify-generated
```

API-facing frontend code should use the generated schema types instead of
handwritten request or response shapes.
