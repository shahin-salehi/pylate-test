from typing import override
from pylate import indexes, retrieve, models
import parse 

#db
chunks = parse.parse_pdf("SEB_Annual_Report_2024_ENG.pdf")
documents_ids = list(range(len(chunks))) 

#model 
model = models.ColBERT(model_name_or_path="lightonai/GTE-ModernColBERT-v1")

# Step 1: Initialize the ColBERT retriever

index = indexes.Voyager(index_folder="pylate-index", index_name="index", override=False)
retriever = retrieve.ColBERT(index=index)


# begin loop
while True: 

    # user input
    user_input = input("query: ")
    if user_input == "exit":
        break
    # Step 2: Encode the queries
    queries_embeddings = model.encode(
        [user_input],
        batch_size=32,
        is_query=True,  
        show_progress_bar=True,
    )

    # Step 3: Retrieve top-k documents
    scores = retriever.retrieve(
        queries_embeddings=queries_embeddings, 
        k=10,  # Retrieve the top 10 matches for each query

    )

    # get from db
    for i, query_result in enumerate(scores):
        for result in query_result:
            print(f"\n--- Score {result["score"]} ---\n{chunks[result["id"]]}\n")




