from unstructured.partition.pdf import partition_pdf
from unstructured.chunking.title import chunk_by_title
from typing import List
import json 

from unstructured_inference.logger import logging 

# logging
log = logging.getLogger(__name__)

def parse_pdf(path: str, category: str) -> tuple[list[dict], bool]:
    try:
        elements = partition_pdf(filename=path, infer_table_structure=True, strategy='hi_res')
        chunks = chunk_by_title(elements)
        tables = [el for el in elements if el.category == "Table"]        

        # add paragraphs
        payloads = []
        for chunk in chunks:
            payloads.append({"filename": chunk.metadata.filename, "pageNumber": chunk.metadata.filename, "category": category, "content": chunk.text, "table": False })

        for table in tables:
            payloads.append({"filename": table.metadata.filename, "pageNumber": table.metadata.filename, "category": category, "content": table.metadata.text_as_html, "table": True })
    
    except Exception as e:
        log.error(f"failed to parse pdf: {path}: {e}")
        return [], False

    return payloads, True 


def main():
    data, ok = parse_pdf("Q125_Quarterly_report.pdf", "test")
    if ok:
        with open("tmp1.json", "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print("Saved payloads to tmp.json")
    else:
        print("Failed to parse PDF.")



main()
