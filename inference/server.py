import logging
import sys
import grpc 
import embed_pb2, embed_pb2_grpc 
import spacy


from concurrent import futures
from internal.colbert.embedder import Embedder
from internal.db.database import Database
from internal.parser.parse import Parse


logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
        )

logger = logging.getLogger(__name__)

# Load the English language model
nlp = spacy.load("en_core_web_sm")


def highlight(paragraph: str, words: list[str]):  
    done = []
    original_paragraph = paragraph  # Keep the original for reference
    highlighted = original_paragraph  # Start with the original paragraph
    # total hack
    searchTerm = original_paragraph[:50]
    
    #filter
    words = [word for word in words if not nlp.vocab[word].is_stop]
    
    for word in words:
        if word in done or not str.isalpha(word):
            continue
            
        i = 0
        while i < len(highlighted):
            # Check if we found the word (case-insensitive) and it's a whole word
            if (highlighted[i:i+len(word)].lower() == word.lower() and 
                (i == 0 or not highlighted[i-1].isalpha()) and 
                (i+len(word) >= len(highlighted) or not highlighted[i+len(word)].isalpha())):
                
                # Wrap the match in span tags
                highlighted = (
                    highlighted[:i] +
                    '<span class="bg-yellow-200 text-black px-1 rounded">' +
                    highlighted[i:i+len(word)] +
                    '</span>' +
                    highlighted[i+len(word):]
                )
                # Skip ahead to avoid overlapping matches
                i += len(word) + len('<span class="bg-yellow-200 text-black px-1 rounded"></span>')
            else:
                i += 1
                
        done.append(word)
    
    return highlighted, searchTerm


class ColBERTEmbedder(embed_pb2_grpc.EmbedderServicer):
    def __init__(self, embedderObject, db):
        self.embedder = embedderObject
        self.db = db
        

    def Embed(self, request, context):

        l, query_embeddings = self.db.search(group=request.group, query=request.text, embedder=self.embedder, category=request.category)

        matches = []
        for row in l:
            doc_tokens, doc_embeds = self.embedder.embed_with_tokens(row[3])
            top_indicies = self.embedder.match(query_embeddings, doc_embeds)
            top_words = [doc_tokens[i] for i in top_indicies]
            #print("top words:", top_words)
            out, st = highlight(row[3], top_words)
            title = f"{row[2]}: {row[0]}"
            matches.append(embed_pb2.Match(
                filename=row[0],
                page_number=row[1],
                title=title,
                category=row[2],
                content= out,
                html=row[4] or "",
                score=row[5],
                meta=st,
                file_url=row[6]
                
                ))

        return embed_pb2.EmbedResponse(result=matches)

def serve():
    # init embedder obj
    embedder = Embedder("model/model.onnx")

    # init database
    db = Database("postgres://admin:password@localhost:9876/documents")
    if db.ping():
        logger.info("ping test successfull")
    else:
        logger.error("ping test failed, shutting down.")
        sys.exit(1)

    # workers python threads not true parallel
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    embed_pb2_grpc.add_EmbedderServicer_to_server(
        ColBERTEmbedder(embedder, db),
        server,
    )



    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server started on port 50051.")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()


