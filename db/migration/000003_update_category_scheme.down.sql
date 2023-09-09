ALTER TABLE category
DROP CONSTRAINT fk_publisher_id;

ALTER TABLE category
DROP COLUMN publisher_id;

