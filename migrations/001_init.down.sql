-- migrations/001_init.down.sql
-- Rollback initial migration by dropping users, accounts, and sessions tables

DROP TABLE sessions;
DROP TABLE accounts;
DROP TABLE users;