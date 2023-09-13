CREATE TABLE
  "blocks" (
    "consensus" bool NOT NULL,
    "difficulty" numeric(50, 0),
    "gas_limit" numeric(100, 0) NOT NULL,
    "gas_used" numeric(100, 0) NOT NULL,
    "hash" bytea NOT NULL,
    "miner_hash" bytea NOT NULL,
    "nonce" bytea NOT NULL,
    "number" int8 NOT NULL,
    "parent_hash" bytea NOT NULL,
    "size" int4,
    "timestamp" timestamp(6) NOT NULL,
    "created_at" timestamp(6) NOT NULL,
    "updated_at" timestamp(6) NOT NULL,
    "refetch_needed" bool DEFAULT false,
    "base_fee_per_gas" numeric(100, 0),
    "is_empty" bool,
    CONSTRAINT "blocks_pkey" PRIMARY KEY ("hash")
  );

ALTER TABLE "blocks" OWNER TO "postgres";

CREATE INDEX "blocks_consensus_index" ON "blocks" USING btree (
  "consensus" "pg_catalog"."bool_ops" ASC NULLS LAST
);

CREATE INDEX "blocks_date" ON "blocks" USING btree (
  date ("timestamp") "pg_catalog"."date_ops" ASC NULLS LAST,
  "number" "pg_catalog"."int8_ops" ASC NULLS LAST
);

CREATE INDEX "blocks_created_at_index" ON "blocks" USING btree (
  "created_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

CREATE INDEX "blocks_is_empty_index" ON "blocks" USING btree ("is_empty" "pg_catalog"."bool_ops" ASC NULLS LAST);

CREATE INDEX "blocks_miner_hash_index" ON "blocks" USING btree (
  "miner_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "blocks_miner_hash_number_index" ON "blocks" USING btree (
  "miner_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "number" "pg_catalog"."int8_ops" ASC NULLS LAST
);

CREATE INDEX "blocks_number_index" ON "blocks" USING btree ("number" "pg_catalog"."int8_ops" ASC NULLS LAST);

CREATE INDEX "blocks_timestamp_index" ON "blocks" USING btree (
  "timestamp" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

CREATE INDEX "empty_consensus_blocks" ON "blocks" USING btree (
  "consensus" "pg_catalog"."bool_ops" ASC NULLS LAST
)
WHERE
  is_empty IS NULL;

CREATE UNIQUE INDEX "one_consensus_block_at_height" ON "blocks" USING btree ("number" "pg_catalog"."int8_ops" ASC NULLS LAST)
WHERE
  consensus;

CREATE UNIQUE INDEX "one_consensus_child_per_parent" ON "blocks" USING btree (
  "parent_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
)
WHERE
  consensus;