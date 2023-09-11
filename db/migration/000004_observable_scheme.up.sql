CREATE TABLE IF NOT EXISTS ip(
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) NOT NULL,
    hostname TEXT,
    city     TEXT,
    region   TEXT,
    country  TEXT,
    loc      TEXT,
    org      TEXT,
    postal   VARCHAR(10),
    timezone TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS observable (
    id SERIAL PRIMARY KEY,
    account_id INT DEFAULT NULL,
    ip_id      INT NOT NULL,
    browser   TEXT,
    os        TEXT,
    os_version TEXT,
    device    TEXT,
    path      TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE observable
ADD CONSTRAINT fk_account_id
FOREIGN KEY (account_id)
REFERENCES account(id)
ON DELETE SET NULL;

ALTER TABLE observable
ADD CONSTRAINT fk_ip_id
FOREIGN KEY (ip_id)
REFERENCES ip(id)
ON DELETE SET NULL;