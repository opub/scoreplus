arg0:  C:\Users\msn\AppData\Local\Temp\go-build827651450\b001\exe\build.exe
-- game

DROP TABLE IF EXISTS game CASCADE;
CREATE TABLE game
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	sport integer,
	hometeam integer,
	awayteam integer,
	homescore integer,
	awayscore integer,
	start timestamp with time zone,
	final boolean,
	venue integer,
	notes integer[]
)
WITH (
	OIDS = FALSE
);
ALTER TABLE game OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE game_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE game_id_seq TO scoreplus_reader;

-- member

DROP TABLE IF EXISTS member CASCADE;
CREATE TABLE member
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	handle text,
	email text,
	firstname text,
	lastname text,
	verified boolean,
	enabled boolean,
	lastactive timestamp with time zone,
	teams integer[],
	follows integer[],
	followers integer[]
)
WITH (
	OIDS = FALSE
);
ALTER TABLE member OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE member_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE member_id_seq TO scoreplus_reader;

-- note

DROP TABLE IF EXISTS note CASCADE;
CREATE TABLE note
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	message text
)
WITH (
	OIDS = FALSE
);
ALTER TABLE note OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE note_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE note_id_seq TO scoreplus_reader;

-- sport

DROP TABLE IF EXISTS sport CASCADE;
CREATE TABLE sport
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	name text
)
WITH (
	OIDS = FALSE
);
ALTER TABLE sport OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE sport_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE sport_id_seq TO scoreplus_reader;

-- team

DROP TABLE IF EXISTS team CASCADE;
CREATE TABLE team
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	name text,
	sport integer,
	venue integer,
	mascot text,
	games integer[]
)
WITH (
	OIDS = FALSE
);
ALTER TABLE team OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE team_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE team_id_seq TO scoreplus_reader;

-- venue

DROP TABLE IF EXISTS venue CASCADE;
CREATE TABLE venue
(
	id serial PRIMARY KEY,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone NOT NULL DEFAULT now(),
	modifiedby integer NOT NULL,
	name text,
	address text,
	coordinates text
)
WITH (
	OIDS = FALSE
);
ALTER TABLE venue OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE venue_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE venue_id_seq TO scoreplus_reader;

