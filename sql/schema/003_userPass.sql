-- +goose Up
ALTER TABLE users 
ADD COLUMN hash TEXT;

-- +goose Down
ALTER TABLE users 
DROP COLUMN hash;
