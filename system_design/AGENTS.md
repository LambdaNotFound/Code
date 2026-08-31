# system_design/

Written system-design solutions and reusable interview notes, in Markdown. Not part of the Go module, not tested, not linked from README.md's problem index.

Two kinds of file here:
- **Worked problems**, each a full mock-interview writeup (functional/non-functional requirements, data model, API, deep dive): `Andruil_URL_Shortener.md`, `Rippling_Autopay_System.md`, `Rippling_Hotel_Booking_System.md`.
- **Reusable concept notes**, referenced across multiple problems rather than tied to one: `consistency_model.md` (quorum/linearizability latency trade-offs), `database_selection_framework.md` (a Postgres-first decision framework — replicate before sharding, shard before switching to Cassandra/DynamoDB), `shards.md` (cluster/shard/replica topology basics), `feed_generation.md` and `numeric_aggregation.md` (ASCII-diagrammed design patterns, each covering a "pair" of related problems).

Company names in problem filenames (Andruil, Rippling) mark which mock-interview context the problem came from, not an affiliation requirement for reading it. When adding a new worked problem, follow the FR/NFR-first structure the existing files use.
