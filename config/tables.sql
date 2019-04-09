-- game

DROP TABLE IF EXISTS game CASCADE;
CREATE TABLE game
(
	id serial NOT NULL PRIMARY KEY,
	sport integer,
	hometeam integer,
	awayteam integer,
	homescore integer,
	awayscore integer,
	start timestamp with time zone,
	final boolean,
	venue integer,
	date timestamp with time zone,
	notes integer[],
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE game OWNER TO scoreplus_owner;

-- member

DROP TABLE IF EXISTS member CASCADE;
CREATE TABLE member
(
	id serial NOT NULL PRIMARY KEY,
	handle text,
	email text,
	firstname text,
	lastname text,
	verified boolean,
	enabled boolean,
	lastactive timestamp with time zone,
	teams integer[],
	follows integer[],
	followers integer[],
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE member OWNER TO scoreplus_owner;

-- note

DROP TABLE IF EXISTS note CASCADE;
CREATE TABLE note
(
	id serial NOT NULL PRIMARY KEY,
	message text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE note OWNER TO scoreplus_owner;

-- sport

DROP TABLE IF EXISTS sport CASCADE;
CREATE TABLE sport
(
	id serial NOT NULL PRIMARY KEY,
	name text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE sport OWNER TO scoreplus_owner;

-- team

DROP TABLE IF EXISTS team CASCADE;
CREATE TABLE team
(
	id serial NOT NULL PRIMARY KEY,
	name text,
	sport integer,
	venue integer,
	mascot text,
	games integer[],
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE team OWNER TO scoreplus_owner;

-- venue

DROP TABLE IF EXISTS venue CASCADE;
CREATE TABLE venue
(
	id serial NOT NULL PRIMARY KEY,
	name text,
	address text,
	coordinates text,
	created timestamp with time zone NOT NULL DEFAULT now(),
	createdby integer NOT NULL,
	modified timestamp with time zone,
	modifiedby integer
)
WITH (
	OIDS = FALSE
);
ALTER TABLE venue OWNER TO scoreplus_owner;

