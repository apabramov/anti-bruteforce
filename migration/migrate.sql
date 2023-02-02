-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklist (
  subnet INET NOT NULL PRIMARY KEY
);

CREATE TABLE whitelist (
  subnet INET NOT NULL PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blacklist;
DROP TABLE IF EXISTS whitelist;
-- +goose StatementEnd
