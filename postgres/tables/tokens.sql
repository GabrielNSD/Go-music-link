BEGIN TRANSACTION;
CREATE TABLE tokens (
  service_name VARCHAR(56) PRIMARY KEY,
  access_token VARCHAR(256) NOT NULL,
  token_type VARCHAR(56),
  scope VARCHAR(256),
  expiration TIMESTAMP,
  refresh_token VARCHAR(256)
);
COMMIT;
