CREATE TABLE teams (
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name TEXT NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    merged_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE pr_reviewers (
    pr_id VARCHAR(255) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id VARCHAR(255) NOT NULL REFERENCES users(id),
    PRIMARY KEY (pr_id, reviewer_id)
);

CREATE INDEX idx_users_team ON users(team_name);
CREATE INDEX idx_pr_author ON pull_requests(author_id);
CREATE INDEX idx_reviews_reviewer ON pr_reviewers(reviewer_id);