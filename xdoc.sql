CREATE SEQUENCE doc_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

CREATE TABLE "public"."doc" (
    "id" bigint DEFAULT nextval('doc_id_seq') NOT NULL,
    "filename" text NOT NULL,
    "username" text NOT NULL,
    "tags" text NOT NULL,
    "data" text NOT NULL,
    "date" integer NOT NULL,
    "time" integer NOT NULL,
    "type" text DEFAULT 'internal' NOT NULL,
    "protect" text DEFAULT 'unprotect' NOT NULL,
    CONSTRAINT "doc_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "doc_filename" ON "public"."doc" USING btree ("filename");


CREATE TABLE "public"."tags" (
    "tagname" text NOT NULL,
    CONSTRAINT "tags_tagname" PRIMARY KEY ("tagname")
) WITH (oids = false);


CREATE TABLE "public"."users" (
    "username" text NOT NULL,
    "canedit" integer NOT NULL,
    "canadmin" integer NOT NULL,
    CONSTRAINT "user_username" PRIMARY KEY ("username")
) WITH (oids = false);
