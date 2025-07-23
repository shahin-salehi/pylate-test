-- activate extension
CREATE EXTENSION IF NOT EXISTS "vector";


-- group
CREATE TABLE IF NOT EXISTS groups (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
    );

-- user
CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
    );

-- user to group
CREATE TABLE IF NOT EXISTS user_to_group(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE
    );

-- pdf
CREATE TABLE IF NOT EXISTS pdfs(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    file_url TEXT NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS pdf_chunks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pdf_id BIGINT REFERENCES pdfs(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    category TEXT,
    content TEXT NOT NULL,
    embeddings vector(128)[] NOT NULL,
    is_table BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS pdf_table_html (
    chunk_id BIGINT PRIMARY KEY REFERENCES pdf_chunks(id) ON DELETE CASCADE,
    html TEXT NOT NULL
);

-- write call that inserts aiops by default at group 1 
