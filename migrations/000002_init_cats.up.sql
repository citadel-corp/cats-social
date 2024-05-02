DROP TYPE IF EXISTS cats_race;
CREATE TYPE cats_race AS ENUM('Persian', 'Maine Coon', 'Siamese', 'Ragdoll',
'Bengal', 'Sphynx', 'British Shorthair', 'Abyssinian', 'Scottish Fold', 'Birman');

DROP TYPE IF EXISTS cats_sex;
CREATE TYPE cats_sex AS ENUM('female', 'male');

CREATE TABLE IF NOT EXISTS
cats (
    id CHAR(16) PRIMARY KEY,
    user_id CHAR(16) NOT NULL,
    name VARCHAR NOT NULL,
	race cats_race NOT NULL,
	sex cats_sex NOT NULL,
	age_in_month INT NOT NULL check (age_in_month between 1 and 120082),
	description TEXT NOT NULL,
	has_matched BOOLEAN NOT NULL DEFAULT FALSE,
	image_urls TEXT[],
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE cats
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS cats_name
	ON cats (name);
CREATE INDEX IF NOT EXISTS cats_user_id
	ON cats USING HASH (user_id);
CREATE INDEX IF NOT EXISTS cats_race_idx
	ON cats(race);
CREATE INDEX IF NOT EXISTS cats_sex_idx
	ON cats(sex);
CREATE INDEX IF NOT EXISTS cats_created_at_desc
	ON cats(created_at DESC);
CREATE INDEX IF NOT EXISTS cats_age
	ON cats(age_in_month);