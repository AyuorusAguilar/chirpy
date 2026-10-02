-- +goose Up
ALTER TABLE users 
ADD hash TEXT DEFAULT 'unset' NOT NULL;

-- +goose Down
ALTER TABLE users 
DROP COLUMN hash;
