DROP TABLE roasters CASCADE;
DROP TABLE roast_sessions CASCADE;
DROP TABLE session_measurements CASCADE;

CREATE TABLE roasters (
  id BIGSERIAL PRIMARY KEY,
  serial_number VARCHAR(36) UNIQUE NOT NULL
);

CREATE TABLE roast_sessions (
  id BIGSERIAL PRIMARY KEY,
  roaster_id INTEGER REFERENCES roasters(id),
  user_id INTEGER REFERENCES users(id),
  roast_date timestamp NOT NULL
);

CREATE TABLE session_measurements (
  session_id INTEGER PRIMARY KEY REFERENCES roast_sessions(id),
  suhu DOUBLE PRECISION
);

INSERT INTO roasters (serial_number) VALUES ('f25a87c9-4613-49da-a61d-b16f352441ad')
