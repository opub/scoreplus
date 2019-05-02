-- game

DROP TABLE IF EXISTS game CASCADE;
CREATE TABLE game
(
	id serial PRIMARY KEY
	, created timestamp with time zone NOT NULL DEFAULT now()
	, createdby integer NOT NULL DEFAULT 0
	, modified timestamp with time zone
	, modifiedby integer NOT NULL DEFAULT 0
	, sport text
	, hometeam integer NOT NULL DEFAULT 0
	, awayteam integer NOT NULL DEFAULT 0
	, homescore integer
	, awayscore integer
	, start timestamp with time zone
	, final boolean
	, venue integer NOT NULL DEFAULT 0
	, notes integer[] DEFAULT '{}'
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
	id serial PRIMARY KEY
	, created timestamp with time zone NOT NULL DEFAULT now()
	, createdby integer NOT NULL DEFAULT 0
	, modified timestamp with time zone
	, modifiedby integer NOT NULL DEFAULT 0
	, handle text NOT NULL UNIQUE
	, email email NOT NULL UNIQUE
	, firstname text
	, lastname text
	, provider text
	, providerid text
	, enabled boolean
	, lastactive timestamp with time zone
	, teams integer[] DEFAULT '{}'
	, follows integer[] DEFAULT '{}'
	, followers integer[] DEFAULT '{}'
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
	id serial PRIMARY KEY
	, created timestamp with time zone NOT NULL DEFAULT now()
	, createdby integer NOT NULL DEFAULT 0
	, modified timestamp with time zone
	, modifiedby integer NOT NULL DEFAULT 0
	, message text
)
WITH (
	OIDS = FALSE
);
ALTER TABLE note OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE note_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE note_id_seq TO scoreplus_reader;

-- team

DROP TABLE IF EXISTS team CASCADE;
CREATE TABLE team
(
	id serial PRIMARY KEY
	, created timestamp with time zone NOT NULL DEFAULT now()
	, createdby integer NOT NULL DEFAULT 0
	, modified timestamp with time zone
	, modifiedby integer NOT NULL DEFAULT 0
	, name text
	, sport text
	, venue integer NOT NULL DEFAULT 0
	, mascot text
	, games integer[] DEFAULT '{}'
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
	id serial PRIMARY KEY
	, created timestamp with time zone NOT NULL DEFAULT now()
	, createdby integer NOT NULL DEFAULT 0
	, modified timestamp with time zone
	, modifiedby integer NOT NULL DEFAULT 0
	, name text
	, address text
	, coordinates text
)
WITH (
	OIDS = FALSE
);
ALTER TABLE venue OWNER TO scoreplus_owner;
GRANT SELECT, USAGE ON SEQUENCE venue_id_seq TO scoreplus_writer;
GRANT SELECT ON SEQUENCE venue_id_seq TO scoreplus_reader;

