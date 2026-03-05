# Adapters Layer Rules

Translates between the outside world and the application core.
**No business logic. No domain objects returned to callers.**

## HTTP Handler Pattern (`in/http/`)

```go
func (h *[X]Handler) [Action](c echo.Context) error {
    // 1. Parse & validate request shape (not business rules)
    var req [Action]Request
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, errorResponse(err))
    }

    // 2. Parse IDs — return 400 on parse failure
    id, err := uuid.Parse(req.ID)
    if err != nil {
        return c.JSON(http.StatusBadRequest, errorResponse("invalid id format"))
    }

    // 3. Build command/query and call handler
    cmd := commands.[Action]Command{ID: id, ...}
    result, err := h.handler.Handle(c.Request().Context(), cmd)
    if err != nil {
        // 422 for domain/business errors, 500 for infrastructure
        return c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
    }

    // 4. Return JSON
    return c.JSON(http.StatusOK, result)
}

func errorResponse(err interface{}) map[string]string {
    return map[string]string{"error": fmt.Sprint(err)}
}
```

## HTTP Status Code Mapping

| Error type | Status |
|-----------|--------|
| Invalid request format, bad UUID | 400 Bad Request |
| Business rule violation (domain error) | 422 Unprocessable Entity |
| Resource not found | 404 Not Found |
| Unauthorized | 401 Unauthorized |
| Infrastructure error (DB, timeout) | 500 Internal Server Error |

## Postgres Repository Pattern (`out/postgres/`)

```go
func (r *Repository) Get(ctx context.Context, tx ports.Tx, id uuid.UUID) (*model.Aggregate, error) {
    query, args, _ := squirrel.
        Select("id", "tenant_id", "field1", "field2").
        From("table_name").
        Where(squirrel.Eq{"id": id, "tenant_id": tenantFromCtx(ctx)}). // ALWAYS include tenant_id
        PlaceholderFormat(squirrel.Dollar).
        ToSql()

    row := r.pool.QueryRow(ctx, query, args...)

    var id uuid.UUID
    var tenantID uuid.UUID
    // ... scan fields

    return model.RestoreAggregate(id, tenantID, ...), nil
}
```

## Rules

- **Always** include `tenant_id` filter in every query — multi-tenancy is not optional
- **Never** use raw SQL string concatenation — always use squirrel
- **Never** return domain objects from HTTP handlers — use response structs
- **Use** `Restore*()` (not `New*()`) when loading aggregates from DB — no validation, no events
- **Map** infrastructure errors (pgx.ErrNoRows → 404) at the adapter boundary
- **Register** all routes in the handler's `Register(e *echo.Echo)` method
