PRAGMA foreign_keys = ON;

CREATE TABLE reservation
(
  rid              integer PRIMARY KEY AUTOINCREMENT,
  seat_count       integer DEFAULT 1,
  start_time       integer,
  customer_id      integer,
  reservation_name text,
  phone            text,
  comments         text,
  created          integer,
  last_updated     integer,
  FOREIGN KEY (rid) REFERENCES customer
);

CREATE TABLE customer
(
  cid          integer PRIMARY KEY AUTOINCREMENT,
  first_name   text NOT NULL,
  last_name    text NOT NULL,
  email        text NOT NULL,
  phone        text,
  created      integer,
  last_updated integer
);