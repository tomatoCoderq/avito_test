-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- teams
CREATE TABLE IF NOT EXISTS teams (
    id VARCHAR(255) PRIMARY KEY DEFAULT uuid_generate_v4()::text,
    name VARCHAR(255) NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS idx_teams_name ON teams(name);

-- teams <-> users (many-to-many)
CREATE TABLE IF NOT EXISTS team_users (
    team_id VARCHAR(255) NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, user_id)
);
-- обратный индекс для быстрого поиска команд пользователя
CREATE INDEX IF NOT EXISTS idx_team_users_user_id ON team_users(user_id);

-- prs
CREATE TABLE IF NOT EXISTS prs (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    need_more_reviewers BOOLEAN NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS idx_prs_author_id ON prs(author_id);
CREATE INDEX IF NOT EXISTS idx_prs_status ON prs(status);

-- prs <-> users (reviewers many-to-many)
CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id   VARCHAR(255) NOT NULL REFERENCES prs(id)   ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (pr_id, user_id)
);
-- быстрый поиск "во что ревьюер добавлен"
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user_id ON pr_reviewers(user_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS team_users;
DROP TABLE IF EXISTS prs;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
