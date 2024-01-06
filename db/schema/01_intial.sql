CREATE SCHEMA IF NOT EXISTS meetup;
-- creat tables for our meetup data
/*
CREATE table if not exists groups (
		id varchar PRIMARY KEY,
		name varchar NOT NULL,
		created_at timetz NOT NULL DEFAULT get_current_time(),
		updated_at timetz NOT NULL DEFAULT get_current_time(),
);
*/
CREATE table if not exists events (
		id varchar PRIMARY KEY,
		title varchar NOT NULL,
		group_id varchar,
		group_name varchar NOT NULL,
		date timestamp,
		going int,	
		waiting int,
		created_at timetz NOT NULL DEFAULT get_current_time(),
		updated_at timetz NOT NULL DEFAULT get_current_time(),
);
