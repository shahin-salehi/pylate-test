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


    def pdf(self, group:int, path: str, file_url:str, filename:str, category: str) -> tuple[dict, bool]:
        try:
            elements = partition_pdf(filename=path, infer_table_structure=True, strategy='hi_res')
            chunks = chunk_by_title(elements)
            tables = [el for el in elements if el.category == "Table"]

            d = {
                "owner": group,
                "filename": filename, 
                "file_url": file_url,
                "chunks": []
            }

            for chunk in chunks:
                embedding = self.embedder.Embed(chunk.text)
                p = {
                    "page_number": chunk.metadata.page_number,
                    "category": category,
                    "content": chunk.text,
                    "embedding": embedding,
                    "is_table": False
                }
                d["chunks"].append(p)

            for table in tables:
                if not table.text:
                    continue  # skip empty tables

                embedding = self.embedder.Embed(table.text)
                html = table.metadata.text_as_html if hasattr(table.metadata, "text_as_html") else ""
                p = {
                    "page_number": table.metadata.page_number,
                    "category": category,
                    "content": table.text,
                    "embedding": embedding,
                    "is_table": True,
                    "html": html
                }
                d["chunks"].append(p)

        except Exception as e:
            logger.error(f"failed to parse pdf: {path}: {e}")
            return {}, False

        return d, True

