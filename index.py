from pylate import indexes, models, retrieve
import parse

chunks = parse.parse_pdf("SEB_Annual_Report_2024_ENG.pdf")


# load model


# Step 1: Load the ColBERT model
model = models.ColBERT(
    model_name_or_path="lightonai/GTE-ModernColBERT-v1",
)

# Step 2: Initialize the Voyager index
index = indexes.Voyager(
    index_folder="pylate-index",
    index_name="index",
    override=True,  # This overwrites the existing index if any
)

# Step 3: Encode the documents
documents_ids = list(range(len(chunks))) 
documents = chunks  

documents_embeddings = model.encode(
    documents,
    batch_size=32,
    is_query=False,  # Ensure that it is set to False to indicate that these are documents, not queries:
    show_progress_bar=True,
)

# Step 4: Add document embeddings to the index by providing embeddings and corresponding ids
index.add_documents(
    documents_ids=documents_ids,
    documents_embeddings=documents_embeddings,
)
