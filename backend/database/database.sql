CREATE DATABASE test_assignment;
-- add uuid generate extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE files ADD COLUMN signature TEXT;
ALTER TABLE files ADD COLUMN signature_date TIMESTAMP;