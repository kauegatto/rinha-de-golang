CREATE TABLE
  clients (
    id INT PRIMARY KEY,
    account_limit INTEGER NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0
  );

CREATE TABLE
  transactions (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    operation CHAR(1) NOT NULL,
    description VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_transactions_client_id FOREIGN KEY (client_id) REFERENCES clients (id)
  );

DO
  $$ BEGIN INSERT INTO clients (id, account_limit) 
VALUES 
  (1, 100000), 
  (2, 80000), 
  (3, 1000000), 
  (4, 10000000), 
  (5, 500000);
END;
$$