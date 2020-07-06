CREATE DATABASE calendar;
\connect calendar;
CREATE SCHEMA calendar;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE calendar.event (
    id varchar NOT NULL DEFAULT uuid_generate_v1(),
    title varchar NOT NULL,
    notice varchar null,
    deleted boolean NULL DEFAULT false,
    date_start timestamp NULL,
    date_complete timestamp NULL,
    creator_id int null,
    editor_id int null,
    date_created timestamp null,
    date_edited timestamp null
);
CREATE UNIQUE INDEX event_id_idx ON calendar.event (id);
