import psycopg
import logging

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
            logger.warning(f"Database ping caused exception: {e}")
            return False

    def insert(self, pdf_chunk):
        print("implement me please")
        

