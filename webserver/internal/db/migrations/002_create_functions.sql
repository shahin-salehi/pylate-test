-- https://github.com/pgvector/pgvector-python/blob/master/examples/colbert/exact.py



CREATE OR REPLACE FUNCTION max_sim(document vector[], query vector[]) RETURNS double precision AS $$
    WITH queries AS (
        SELECT row_number() OVER () AS query_number, * FROM (SELECT unnest(query) AS query)
    ),
    documents AS (
        SELECT unnest(document) AS document
    ),
    similarities AS (
        SELECT query_number, 1 - (document <=> query) AS similarity FROM queries CROSS JOIN documents
    ),
    max_similarities AS (
        SELECT MAX(similarity) AS max_similarity FROM similarities GROUP BY query_number
    )
    SELECT SUM(max_similarity) FROM max_similarities
$$ LANGUAGE SQL;

-------------- read_files


CREATE OR REPLACE FUNCTION read_files(
    group_id BIGINT
)
RETURNS TABLE(
    pdf_id BIGINT,
    filename TEXT,
    uploaded_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.filename, p.uploaded_at
    FROM pdfs p
    WHERE p.owner = group_id
    ORDER BY p.uploaded_at DESC;
END;
$$ LANGUAGE plpgsql;


--- get user 
CREATE OR REPLACE FUNCTION get_user_by_email(_email TEXT)
RETURNS TABLE (
    id BIGINT,
    password_hash TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT users.id, users.password_hash
    FROM users
    WHERE users.email = _email;
END;
$$ LANGUAGE plpgsql;


--- register user
CREATE OR REPLACE FUNCTION register_user(_username TEXT, _email TEXT, _password_hash TEXT)
RETURNS BIGINT AS $$
DECLARE
    _id BIGINT;
BEGIN
    INSERT INTO users (username, email, password_hash)
    VALUES (_username, _email, _password_hash)
    RETURNING id INTO _id;

    RETURN _id;
END;
$$ LANGUAGE plpgsql;


