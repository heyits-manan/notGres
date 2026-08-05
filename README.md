# notGres

*a database server that's definitely not Postgres*

notGres is a from-scratch, learning-first database server written in Go. It speaks (or will speak) the PostgreSQL wire protocol, because why design your own protocol when the world already has a perfectly good one? The end goal is a SQL database you can connect to with `psql` and actually run queries on. Emphasis on "end goal" — we're not there yet, pal.

## What is this?

I'm learning Go by building a database. Like, a real one. With a TCP server, a query parser, a storage engine, the whole nine yards. Is this the most efficient way to learn a language? No. Is it the most fun? Absolutely.

The architecture is *heavily inspired* by Postgres because, well, it's the gold standard. The internal package structure mirrors how real databases are organized: a server layer, a parser, a planner/executor, and a storage engine. Right now most of those folders are empty, but hey — you gotta start somewhere.

## Current features

- **Postgres wire protocol startup handshake.** `psql` can connect and see a prompt. We speak the language now.
- **TCP listener** on port `5432` (because we're cosplaying Postgres).
- **Binary protocol primitives** — manual `ReadInt32`/`WriteInt32`/`ReadString`/`WriteString`. No libraries. Raw bytes over the wire.
- **Graceful shutdown.** Hit Ctrl+C and it cleans up like a responsible adult.
- **Signal handling.** SIGTERM, SIGINT — we respect them all.
- **A Makefile** that does `build`, `run`, `test`, and `vet`. Fancy.
- **Zero dependencies.** Pure stdlib. No `pgx`, no `pq`, no nothing.

## What's next (the roadmap)

- **Query loop** — read SQL text from the client, send back a hardcoded result (hello `SELECT 1;`)
- **SQL parser** — lexer, tokenizer, AST. The fun recursive descent stuff
- **Query planner & executor** — turning parsed queries into actual work
- **Storage engine** — pages, heap files, a buffer pool. On-disk persistence!
- **Catalog** — table metadata, schemas, the system tables
- **Transactional support** — ACID, baby. Well, maybe A and D first
- **Indexes** — B+ trees or bust
- **WAL (Write-Ahead Log)** — crash recovery so your data doesn't vanish into the void

## Why "notGres"?

Because it's not Postgres. But it wants to be when it grows up.

## Quick start

```bash
# clone
git clone https://github.com/heyits-manan/notGres.git
cd notGres

# run
make run

# in another terminal — psql can now connect successfully!
psql -h localhost -p 5432
# psql (14.0)
# Type "help" for help.
# (no database)=#

# build
make build

# test
make test
```

## Project structure

```
notGres/
├── main.go              # entrypoint: flags, listener, signal handling
├── internal/
│   ├── protocol/        # Postgres wire protocol: read/write primitives, startup handshake
│   ├── server/          # TCP server, connection handling
│   ├── sql/             # future: lexer, parser, AST
│   ├── executor/        # future: query execution engine
│   ├── storage/         # future: heap files, pages, buffer pool
│   ├── catalog/         # future: table & schema metadata
│   └── types/           # future: column types & value encodings
├── Makefile             # build, run, test, vet
├── go.mod               # module definition (zero deps, baby)
└── data/                # database files (gitignored)
```

## License

MIT — go wild. Build cooler things on top of this. Then tell me about them.
