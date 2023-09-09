ALTER TABLE category
ADD COLUMN publisher_id INT DEFAULT NULL;

ALTER TABLE category
ADD CONSTRAINT fk_publisher_id
FOREIGN KEY (publisher_id)
REFERENCES account(id)
ON DELETE SET NULL;