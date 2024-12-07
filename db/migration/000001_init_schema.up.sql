CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz DEFAULT 'now()'
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "owner_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT 'now()'
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_acount_id" bigint NOT NULL,
  "to_acount_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT 'now()'
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("owner_id");

CREATE INDEX ON "transfers" ("from_acount_id");

CREATE INDEX ON "transfers" ("to_acount_id");

CREATE INDEX ON "transfers" ("from_acount_id", "to_acount_id");

ALTER TABLE "entries" ADD FOREIGN KEY ("owner_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_acount_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_acount_id") REFERENCES "accounts" ("id");
