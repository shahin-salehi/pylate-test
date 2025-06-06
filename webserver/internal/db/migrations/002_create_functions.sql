-- https://github.com/pgvector/pgvector-python/blob/master/examples/colbert/exact.py
-- max sim based on query
CREATE OR REPLACE FUNCTION max_sim_chunks(query vector[])
RETURNS TABLE (
    pdf_id BIGINT,
    score DOUBLE PRECISION
)
AS $$
    WITH queries AS (
        SELECT row_number() OVER () AS query_number, unnest(query) AS qvec
    ),
    similarities AS (
        SELECT
            queries.query_number,
            c.pdf_id,
            1 - (c.embedding <=> queries.qvec) AS similarity
        FROM queries
        JOIN pdf_chunks c ON TRUE
    ),
    max_similarities AS (
        SELECT
            pdf_id,
            query_number,
            MAX(similarity) AS max_similarity
        FROM similarities
        GROUP BY pdf_id, query_number
    ),
    summed_scores AS (
        SELECT
            pdf_id,
            SUM(max_similarity) AS score
        FROM max_similarities
        GROUP BY pdf_id
    )
    SELECT pdf_id, score
    FROM summed_scores
    ORDER BY score DESC;
$$ LANGUAGE SQL STABLE;

-- insert pdf chunk
CREATE OR REPLACE FUNCTION insert_pdf_from_json(p_data JSONB)
RETURNS BIGINT AS $$
DECLARE
    new_pdf_id BIGINT;
    chunk JSONB;
BEGIN
    -- Insert PDF row
    INSERT INTO pdfs (owner, filename)
    VALUES (
        (p_data->>'owner')::BIGINT,
        p_data->>'filename'
    )
    RETURNING id INTO new_pdf_id;

    -- Loop over chunks and insert each one
    FOR chunk IN SELECT * FROM jsonb_array_elements(p_data->'chunks')
    LOOP
        INSERT INTO pdf_chunks (
            pdf_id,
            page_number,
            category,
            content,
            embedding,
            is_table
        ) VALUES (
            new_pdf_id,
            (chunk->>'page_number')::INT,
            chunk->>'category',
            chunk->>'content',
            (chunk->>'embedding')::vector,
            (chunk->>'is_table')::BOOLEAN
        );
    END LOOP;

    RETURN new_pdf_id;
END;
$$ LANGUAGE plpgsql;

