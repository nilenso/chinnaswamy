CREATE TABLE url_mappings (
  short_url text,
  long_url text,
  valid boolean,
  created_at timestamp,
  PRIMARY KEY (short_url)
);
