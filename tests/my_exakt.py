
from colbert.infra import ColBERTConfig
from colbert.modeling.checkpoint import Checkpoint
from pgvector.psycopg import register_vector
import psycopg
import warnings

conn = psycopg.connect(
    host="localhost",
    port=9876,
    user="admin",
    password="password",
    dbname="documents",
    autocommit=True
)

conn.execute('CREATE EXTENSION IF NOT EXISTS vector')
register_vector(conn)

#conn.execute('DROP TABLE IF EXISTS documents')
#conn.execute('CREATE TABLE documents (id bigserial PRIMARY KEY, content text, embeddings vector(128)[])')
conn.execute("""
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
$$ LANGUAGE SQL
""")

warnings.filterwarnings('ignore')  # ignore warnings from colbert

# 🧠 Force CPU and float32 vectors
config = ColBERTConfig(
    doc_maxlen=220,
    query_maxlen=32,
    nbits=32,
)

checkpoint = Checkpoint('colbert-ir/colbertv2.0', colbert_config=config, verbose=0)
"""
input = [
    'The dog is barking',
    'The cat is purring',
    'The bear is growling'
onnx model output not normalized]

doc_embeddings = checkpoint.docFromText(input, keep_dims=False)
for content, embeddings in zip(input, doc_embeddings):
    # 🧠 Move to CPU before converting to NumPy
    embeddings = [e.cpu().numpy() for e in embeddings]
    conn.execute('INSERT INTO documents (content, embeddings) VALUES (%s, %s)', (content, embeddings))

"""
query = 'transition plan'
query_embeddings = [e.cpu().numpy() for e in checkpoint.queryFromText([query])[0]]

print("embeddings:", query_embeddings[0])

result = conn.execute(
    'SELECT content, max_sim(embeddings, %s) AS max_sim FROM pdf_chunks ORDER BY max_sim DESC LIMIT 10',
    (query_embeddings[:3],)
).fetchall()

for row in result:
    print(row)
    
    print("\n\n")


conn.close()

