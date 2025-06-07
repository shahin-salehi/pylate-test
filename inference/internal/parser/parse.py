import logging
from unstructured.partition.pdf import partition_pdf
from unstructured.chunking.title import chunk_by_title
from internal.colbert.embedder import Embedder
# Configure the logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)  # Change to DEBUG for more verbose output


"""
when I first started I had ambitions to type everything
that ambition has since passed and other files will
be what they are.

"""

class Parse:
    def __init__(self, embedder: Embedder): 
        self.embedder = embedder

    def pdf(self, path: str, filename: str, category: str) -> tuple[dict, bool]:
        try:
            elements = partition_pdf(filename=path, infer_table_structure=True, strategy='hi_res')
            # parallelize if possible
            chunks = chunk_by_title(elements)
            tables = [el for el in elements if el.category == "Table"]

            d = {
                    "owner": 1, #this will be send from the requesting webserver
                    "filename": filename,
                    "chunks": []
                }
            for chunk in chunks:
                p = {
                        "page_number": chunk.metadata.page_number,
                        "category": category,
                        "content": chunk.text,
                        "embedding": self.embedder.Embed(chunk.text),
                        "is_table": False
                    }
                d["chunks"].append(p)

            for table in tables:
                p = {
                        "page_number": chunk.metadata.page_number,
                        "category": category,
                        "content": chunk.text,
                        "embedding": self.embedder.Embed(chunk.text),
                        "is_table": True,
                        "html": table.metadata.text_as_html
                    }
                d["chunks"].append(p)

        
        except Exception as e:
            logger.error(f"failed to parse pdf: {path}: {e}")
            return {}, False

        return d, True 

