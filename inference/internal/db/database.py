from typing import Tuple
import psycopg
import numpy as np
import logging
from psycopg.types.json import Json 
from internal.colbert.embedder import Embedder
from pgvector.psycopg import register_vector

# Configure the logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)  # Change to DEBUG for more verbose output


"""
Since the DB connection on the inference side only handles inserts, 
we do not prioritize speed and run it sync.

The webserver will have an async pool of connections.
"""

def make_json_safe(obj):
    if isinstance(obj, np.ndarray):
        return obj.tolist()
    elif isinstance(obj, list):
        return [make_json_safe(x) for x in obj]
    elif isinstance(obj, dict):
        return {k: make_json_safe(v) for k, v in obj.items()}
    return obj


class Database:
    def __init__(self, connection_string):
        try:
            self.conn = psycopg.connect(connection_string)
            # register vector
            register_vector(self.conn)
            logger.info("Database connection established.")
        except Exception as e:
            logger.error(f"Failed to connect to database: {e}")
            raise

    def ping(self):
        try:
            with self.conn.cursor() as c:
                c.execute("SELECT 1")
                return c.fetchone() is not None
                
        except Exception as e:
            logger.error(f"Database ping caused exception: {e}")
            return False

    def insert(self, json_data):
        try:
            with self.conn.cursor() as c:
                c.execute("SELECT insert_pdf_from_json(%s::jsonb)", (Json(make_json_safe(json_data)),))
                row = c.fetchone()
                if row is not None:
                    pdf_id = row[0]
                    self.conn.commit()
                    return pdf_id, True
                else:
                    return -1, False # if this happens everyone will be confused

        except Exception as e:
            logger.error(f"Insert parsed pdf caused exception: {e}")
            self.conn.rollback()
            return e, False


    def insert_pdf(self, data: dict) -> tuple[int, bool]:
        try:
            with self.conn.cursor() as c:
                # Step 1: Insert into pdfs
                c.execute(
                    "INSERT INTO pdfs (owner, filename) VALUES (%s, %s) RETURNING id",
                    (data["owner"], data["filename"])
                )
                pdf_id = c.fetchone()[0]

                # Step 2: Insert each chunk
                for chunk in data["chunks"]:
                    c.execute(
                        """
                        INSERT INTO pdf_chunks (
                            pdf_id, page_number, category, content, embeddings, is_table
                        ) VALUES (%s, %s, %s, %s, %s, %s)
                        RETURNING id
                        """,
                        (
                            pdf_id,
                            chunk["page_number"],
                            chunk["category"],
                            chunk["content"],
                            chunk["embedding"],  # List[List[float]] or List[np.ndarray]
                            chunk["is_table"]
                        )
                    )
                    chunk_id = c.fetchone()[0]

                    # If it's a table and has HTML, insert into pdf_table_html
                    if chunk["is_table"] and "html" in chunk:
                        c.execute(
                            "INSERT INTO pdf_table_html (chunk_id, html) VALUES (%s, %s)",
                            (chunk_id, chunk["html"])
                        )

                self.conn.commit()
                return pdf_id, True

        except Exception as e:
            logger.error(f"Insert failed: {e}")
            self.conn.rollback()
            return -1, False
            
    def search(self, group: int, query: str, embedder: Embedder, top_k: int = 10, category: str = "") -> Tuple[list[dict], list]:
        try:
            query_embeddings = embedder.Embed(query)

            sql = """
                SELECT 
                  pdfs.filename,
                  pdf_chunks.page_number,
                  pdf_chunks.category,
                  pdf_chunks.content,
                  pdf_table_html.html,
                  max_sim(embeddings, %s) AS max_sim,
                  pdfs.file_url
                FROM pdf_chunks
                JOIN pdfs ON pdfs.id = pdf_chunks.pdf_id
                LEFT JOIN pdf_table_html ON pdf_table_html.chunk_id = pdf_chunks.id
                WHERE pdfs.owner = %s
            """
            params = [query_embeddings, group]

            if category:
                sql += " AND pdf_chunks.category = %s"
                params.append(category)

            sql += """
                ORDER BY max_sim DESC
                LIMIT %s
            """
            params.append(top_k)

            result = self.conn.execute(sql, tuple(params)).fetchall()
            return result, query_embeddings

        except Exception as e:
            logger.error(f"Search failed: {e}")
            return [], []
