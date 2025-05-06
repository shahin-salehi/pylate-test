import PyPDF2
from typing import List

CHUNK_SIZE = 300  # words

def extract_pdf_text(pdf_path) -> str:
    text = ""
    with open(pdf_path, "rb") as f:
        reader = PyPDF2.PdfReader(f)
        for page in reader.pages:
            text += page.extract_text() + "\n"
    return text

def split_into_chunks(text, max_words=CHUNK_SIZE) -> List[str]:
    words = text.split()
    chunks = []
    for i in range(0, len(words), max_words):
        chunk = " ".join(words[i:i + max_words])
        chunks.append(chunk)
    return chunks

def parse_pdf(path: str) -> List[str]:
    full_text = extract_pdf_text(path)
    return split_into_chunks(full_text)



""" 
nice print example:

    for i, chunk in enumerate(chunks):
        print(f"\n--- Chunk {i+1} ---\n{chunk}\n")

"""
    

