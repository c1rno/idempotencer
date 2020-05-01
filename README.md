# Idempotencer

With project aims to implement inboxer/outboxer patterns over Kafka
 as upstream, Postgres as persistent storage and ZeroMQ as downstream.

There may be alternatives in the future to change persistent storage.

## Rules
	- Prefer internal `errors.Error` over standart `error`-interface
