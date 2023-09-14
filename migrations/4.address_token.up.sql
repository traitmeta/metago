CREATE TABLE "address_current_token_balances" (
  "id" int8 NOT NULL DEFAULT nextval('address_current_token_balances_id_seq'::regclass),
  "address_hash" bytea NOT NULL,
  "block_number" int8 NOT NULL,
  "token_contract_address_hash" bytea NOT NULL,
  "value" numeric,
  "value_fetched_at" timestamp(6),
  "inserted_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "old_value" numeric,
  "token_id" numeric(78,0),
  "token_type" varchar(255) COLLATE "pg_catalog"."default",
  CONSTRAINT "address_current_token_balances_pkey" PRIMARY KEY ("id")
)
;

CREATE TABLE "address_token_balances" (
  "id" int8 NOT NULL DEFAULT nextval('address_token_balances_id_seq'::regclass),
  "address_hash" bytea NOT NULL,
  "block_number" int8 NOT NULL,
  "token_contract_address_hash" bytea NOT NULL,
  "value" numeric,
  "value_fetched_at" timestamp(6),
  "inserted_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "token_id" numeric(78,0),
  "token_type" varchar(255) COLLATE "pg_catalog"."default",
  CONSTRAINT "address_token_balances_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "address_token_balances_token_contract_address_hash_fkey" FOREIGN KEY ("token_contract_address_hash") REFERENCES "tokens" ("contract_address_hash") ON DELETE NO ACTION ON UPDATE NO ACTION
)
;

CREATE TABLE "addresses" (
  "fetched_coin_balance" numeric(100,0),
  "fetched_coin_balance_block_number" int8,
  "hash" bytea NOT NULL,
  "contract_code" bytea,
  "inserted_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "nonce" int4,
  "decompiled" bool,
  "verified" bool,
  "gas_used" int8,
  "transactions_count" int4,
  "token_transfers_count" int4,
  CONSTRAINT "addresses_pkey" PRIMARY KEY ("hash")
)
;

ALTER TABLE "addresses" 
  OWNER TO "postgres";

CREATE INDEX "addresses_fetched_coin_balance_hash_index" ON "addresses" USING btree (
  "fetched_coin_balance" "pg_catalog"."numeric_ops" DESC NULLS FIRST,
  "hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
) WHERE fetched_coin_balance > 0::numeric;

CREATE INDEX "addresses_fetched_coin_balance_index" ON "addresses" USING btree (
  "fetched_coin_balance" "pg_catalog"."numeric_ops" ASC NULLS LAST
);

CREATE INDEX "addresses_inserted_at_index" ON "addresses" USING btree (
  "inserted_at" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);


CREATE TABLE "tokens" (
  "name" text COLLATE "pg_catalog"."default",
  "symbol" text COLLATE "pg_catalog"."default",
  "total_supply" numeric,
  "decimals" numeric,
  "type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "cataloged" bool DEFAULT false,
  "contract_address_hash" bytea NOT NULL,
  "inserted_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "holder_count" int4,
  "skip_metadata" bool,
  "fiat_value" numeric,
  "circulating_market_cap" numeric,
  "total_supply_updated_at_block" int8,
  "icon_url" varchar(255) COLLATE "pg_catalog"."default",
  "is_verified_via_admin_panel" bool DEFAULT false,
  CONSTRAINT "tokens_pkey" PRIMARY KEY ("contract_address_hash")
)
;

ALTER TABLE "tokens" 
  OWNER TO "postgres";

CREATE UNIQUE INDEX "tokens_contract_address_hash_index" ON "tokens" USING btree (
  "contract_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "tokens_symbol_index" ON "tokens" USING btree (
  "symbol" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE INDEX "tokens_trgm_idx" ON "tokens" USING gin (
  to_tsvector('english'::regconfig, (symbol || ' '::text) || name) "pg_catalog"."tsvector_ops"
);

CREATE INDEX "tokens_type_index" ON "tokens" USING btree (
  "type" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

CREATE TABLE "token_transfers" (
  "transaction_hash" bytea NOT NULL,
  "log_index" int4 NOT NULL,
  "from_address_hash" bytea NOT NULL,
  "to_address_hash" bytea NOT NULL,
  "amount" numeric,
  "token_id" numeric(78,0),
  "token_contract_address_hash" bytea NOT NULL,
  "inserted_at" timestamp(6) NOT NULL,
  "updated_at" timestamp(6) NOT NULL,
  "block_number" int4,
  "block_hash" bytea NOT NULL,
  "amounts" numeric[],
  "token_ids" numeric(78,0)[],
  CONSTRAINT "token_transfers_pkey" PRIMARY KEY ("transaction_hash", "block_hash", "log_index"),
  CONSTRAINT "token_transfers_block_hash_fkey" FOREIGN KEY ("block_hash") REFERENCES "blocks" ("hash") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "token_transfers_transaction_hash_fkey" FOREIGN KEY ("transaction_hash") REFERENCES "transactions" ("hash") ON DELETE CASCADE ON UPDATE NO ACTION
)
;

ALTER TABLE "token_transfers" 
  OWNER TO "postgres";

CREATE INDEX "token_transfers_block_number_DESC_log_index_DESC_index" ON "token_transfers" USING btree (
  "block_number" "pg_catalog"."int4_ops" DESC NULLS FIRST,
  "log_index" "pg_catalog"."int4_ops" DESC NULLS FIRST
);

CREATE INDEX "token_transfers_block_number_index" ON "token_transfers" USING btree (
  "block_number" "pg_catalog"."int4_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_from_address_hash_transaction_hash_index" ON "token_transfers" USING btree (
  "from_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_to_address_hash_transaction_hash_index" ON "token_transfers" USING btree (
  "to_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_token_contract_address_hash_block_number_index" ON "token_transfers" USING btree (
  "token_contract_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "block_number" "pg_catalog"."int4_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_token_contract_address_hash_token_id_DESC_block" ON "token_transfers" USING btree (
  "token_contract_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "token_id" "pg_catalog"."numeric_ops" DESC NULLS FIRST,
  "block_number" "pg_catalog"."int4_ops" DESC NULLS FIRST
);

CREATE INDEX "token_transfers_token_contract_address_hash_transaction_hash_in" ON "token_transfers" USING btree (
  "token_contract_address_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_token_id_index" ON "token_transfers" USING btree (
  "token_id" "pg_catalog"."numeric_ops" ASC NULLS LAST
);

CREATE INDEX "token_transfers_token_ids_index" ON "token_transfers" USING gin (
  "token_ids" "pg_catalog"."array_ops"
);

CREATE INDEX "token_transfers_transaction_hash_log_index_index" ON "token_transfers" USING btree (
  "transaction_hash" "pg_catalog"."bytea_ops" ASC NULLS LAST,
  "log_index" "pg_catalog"."int4_ops" ASC NULLS LAST
);