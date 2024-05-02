DROP TYPE IF EXISTS cat_matches_approval;
CREATE TYPE cat_matches_approval AS ENUM('approved', 'rejected', 'pending');

CREATE TABLE IF NOT EXISTS
cat_matches(
    id CHAR(16) PRIMARY KEY,
    issuer_user_id CHAR(16) NOT NULL,
	issuer_cat_id CHAR(16) NOT NULL,
	matched_user_id CHAR(16) NOT NULL,
	matched_cat_id CHAR(16) NOT NULL,
	message VARCHAR(120) NOT NULL,
	approval_status cat_matches_approval NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE cat_matches
	ADD CONSTRAINT fk_issuer_user_id FOREIGN KEY (issuer_user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE cat_matches
	ADD CONSTRAINT fk_matched_user_id FOREIGN KEY (matched_user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS cat_matches_approval_idx
	ON cat_matches(approval_status);
CREATE INDEX IF NOT EXISTS cat_matches_created_at_desc
	ON cat_matches(created_at DESC);
CREATE INDEX IF NOT EXISTS cat_matches_issuer_user_id
	ON cat_matches USING HASH(issuer_user_id);
CREATE INDEX IF NOT EXISTS cat_matches_matched_user_id
	ON cat_matches USING HASH(matched_user_id);