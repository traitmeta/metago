CREATE TABLE "events" (
  "data" bytea NOT NULL,
  "index" int4 NOT NULL,
  "type" varchar(255) COLLATE "pg_catalog"."default",
  "first_topic" varchar(255) COLLATE "pg_catalog"."default",
  "second_topic" varchar(255) COLLATE "pg_catalog"."default",
  "third_topic" varchar(255) COLLATE "pg_catalog"."default",
  "fourth_topic" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "address_hash" bytea,
  "transaction_hash" bytea NOT NULL,
  "block_hash" bytea NOT NULL,
  "block_number" int4,
  CONSTRAINT "events_pkey" PRIMARY KEY ("transaction_hash", "block_hash", "index"),
  CONSTRAINT "events_block_hash_fkey" FOREIGN KEY ("block_hash") REFERENCES "blocks" ("hash") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "events_transaction_hash_fkey" FOREIGN KEY ("transaction_hash") REFERENCES "transactions" ("hash") ON DELETE CASCADE ON UPDATE NO ACTION
)
;

ALTER TABLE "events" 
  OWNER TO "postgres";

CREATE INDEX "events_address_hash_index" ON "events" USING btree (
  "address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "events_address_hash_transaction_hash_index" ON "events" USING btree (
  "address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "events_block_number_DESC__index_DESC_index" ON "events" USING btree (
  "block_number" "pg_catalog"."int4_ops" DESC NULLS FIRST,
  "index" "pg_catalog"."int4_ops" DESC NULLS FIRST
);

CREATE INDEX "events_first_topic_index" ON "events" USING btree (
  "first_topic" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE INDEX "events_fourth_topic_index" ON "events" USING btree (
  "fourth_topic" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE INDEX "events_index_index" ON "events" USING btree (
  "index" "pg_catalog"."int4_ops" ASC NULLS LAST
);

CREATE INDEX "events_second_topic_index" ON "events" USING btree (
  "second_topic" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE INDEX "events_third_topic_index" ON "events" USING btree (
  "third_topic" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE INDEX "events_transaction_hash_index_index" ON "events" USING btree (
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "index" "pg_catalog"."int4_ops" ASC NULLS LAST
);

CREATE INDEX "events_type_index" ON "events" USING btree (
  "type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);