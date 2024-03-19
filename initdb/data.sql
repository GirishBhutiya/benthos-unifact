CREATE TABLE IF NOT EXISTS "data" (
  "id" bigserial PRIMARY KEY,
  "groupdata" varchar,
  "opcua_path" varchar,
  "historian" varchar
);