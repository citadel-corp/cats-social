CREATE TYPE cats_race AS ENUM('Persian', 'Maine Coon', 'Siamese', 'Ragdoll',
'Bengal', 'Sphynx', 'British Shorthair', 'Abyssinian', 'Scottish Fold', 'Birman');

CREATE TYPE cats_sex AS ENUM('Female', 'Male');

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
    CONSTRAINT fk_user_id
			FOREIGN KEY (user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
			ON UPDATE NO ACTION
);

CREATE INDEX IF NOT EXISTS cats_name
	ON cats (name);
CREATE INDEX IF NOT EXISTS cats_user_id
	ON cats USING HASH (user_id);
CREATE INDEX IF NOT EXISTS cats_race
	ON cats USING gin(race);
CREATE INDEX IF NOT EXISTS cats_sex
	ON cats USING gin(sex);
CREATE INDEX IF NOT EXISTS cats_created_at_desc
	ON cats(created_at DESC);
CREATE INDEX IF NOT EXISTS cats_age
	ON cats(age_in_month);