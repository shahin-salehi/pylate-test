import psycopg
import logging
from psycopg.types.json import Json 
# Configure the logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)  # Change to DEBUG for more verbose output


"""
Since the DB connection on the inference side only handles inserts, 
we do not prioritize speed and run it sync.

The webserver will have an async pool of connections.
"""
class Database:
    def __init__(self, connection_string):
        try:
            self.conn = psycopg.connect(connection_string)
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
                c.execute("SELECT insert_pdf_from_json(%s::jsonb)", (Json(json_data),))
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



        

        

