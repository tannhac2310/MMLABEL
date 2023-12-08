-- add sort column for devices table
ALTER TABLE devices ADD COLUMN sort INTEGER NOT NULL DEFAULT 1;