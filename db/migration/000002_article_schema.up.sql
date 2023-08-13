CREATE TABLE IF NOT EXISTS category(
    id SERIAL PRIMARY KEY,
    title TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS article(
    id SERIAL PRIMARY KEY,
    publisher_id INT NOT NULL,
    title TEXT UNIQUE NOT NULL,
    reading_time INT NOT NULL,
    location TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE article
ADD CONSTRAINT fk_publisher_id
FOREIGN KEY (publisher_id)
REFERENCES account(id)
ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS article_category(
    article_id INT,
    category_id INT,
    PRIMARY KEY (article_id, category_id)
);

ALTER TABLE article_category
ADD CONSTRAINT fk_article_id
FOREIGN KEY (article_id)
REFERENCES article(id)
ON DELETE CASCADE;

ALTER TABLE article_category
ADD CONSTRAINT fk_category_id
FOREIGN KEY (category_id)
REFERENCES category(id)
ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS attachment (
    id SERIAL PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS article_attachment(
    article_id INT,
    attachment_id INT,
    PRIMARY KEY (article_id, attachment_id)
);

ALTER TABLE article_attachment
ADD CONSTRAINT fk_article_id
FOREIGN KEY (article_id)
REFERENCES article(id)
ON DELETE CASCADE;

ALTER TABLE article_attachment
ADD CONSTRAINT fk_attachment_id
FOREIGN KEY (attachment_id)
REFERENCES attachment(id)
ON DELETE CASCADE;