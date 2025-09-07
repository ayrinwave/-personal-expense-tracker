CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     email TEXT UNIQUE NOT NULL,
                                     password TEXT NOT NULL,
                                     created_at TIMESTAMP DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS expenses (
                                        id SERIAL PRIMARY KEY,
                                        user_id INT REFERENCES users(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id),
    amount NUMERIC(10,2) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW()
    );
