-- roles

CREATE ROLE scoreplus_writer
  NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

CREATE ROLE scoreplus_reader
  NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

-- logins

CREATE ROLE scoreplus_owner LOGIN
  NOSUPERUSER INHERIT CREATEDB CREATEROLE NOREPLICATION;
GRANT scoreplus_reader TO scoreplus_owner;
GRANT scoreplus_writer TO scoreplus_owner;

CREATE ROLE scoreplus_user LOGIN
  NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
GRANT scoreplus_reader TO scoreplus_user;
GRANT scoreplus_writer TO scoreplus_user;

CREATE ROLE scoreplus_reader LOGIN
  NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
GRANT scoreplus_reader TO scoreplus_reader;

-- database

CREATE DATABASE scoreplus_dev
  WITH OWNER = scoreplus_owner
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       LC_COLLATE = 'English_United States.1252'
       LC_CTYPE = 'English_United States.1252'
       CONNECTION LIMIT = -1;

-- privileges

ALTER DEFAULT PRIVILEGES 
    GRANT INSERT, SELECT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON TABLES
    TO postgres;

ALTER DEFAULT PRIVILEGES 
    GRANT SELECT ON TABLES
    TO scoreplus_reader;

ALTER DEFAULT PRIVILEGES 
    GRANT INSERT, SELECT, UPDATE, DELETE, TRUNCATE, REFERENCES, TRIGGER ON TABLES
    TO scoreplus_writer;

-- custom types

CREATE EXTENSION citext;

CREATE DOMAIN email AS citext
  CHECK ( value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );