# bookstore
Example project of using wicked-sqlc + wpgx + dcache to build scalable service.

Generate with the wicked sqlc fork:

```bash
make sqlc SQLC=/absolute/path/to/sqlc/bin/sqlc
make sqlc-verify SQLC=/absolute/path/to/sqlc/bin/sqlc
```

The upstream-sync fixtures use a fork built with `make build COMMIT_HASH=v2.4.0-dev`.
Runtime versions remain wpgx v0.3.1 and dcache v0.1.3.

Run tests with Go and Docker available:

```bash
make test
go test -race -count=1 -p 1 -timeout 10m ./pkg/usecases
make lint-fix
```

Tests create their own PostgreSQL and Redis containers on random loopback ports
and remove those exact containers afterward. They never use an existing database
or Redis instance. Replica tests use a second connection configuration for the
same test database; they validate generated routing/API behavior, not replication lag.

Coverage includes cache hits and negative caching, nullable pointer cache keys,
JSON round trips, copyfrom, nil-cache mutations, transaction commit/rollback
invalidation, multiple invalidation targets, and cache/database timeouts.
