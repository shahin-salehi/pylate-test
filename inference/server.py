import logging
import sys
import grpc 
import embed_pb2, embed_pb2_grpc 
import json

from concurrent import futures
from internal.colbert.embedder import Embedder
from internal.db.database import Database
from internal.parser.parse import Parse



logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
        )

logger = logging.getLogger(__name__)

class ColBERTEmbedder(embed_pb2_grpc.EmbedderServicer):
    def __init__(self, embedderObject):
        self.embedder = embedderObject
        

    def Embed(self, request, context):
        outputs = self.embedder.Embed(request.text)
        return embed_pb2.EmbedResponse(embedding=outputs)

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

    # init parser
    parse = Parse(embedder)
    
    
    data, ok = parse.pdf("internal/parser/docs/example_pdf.pdf", "test", "test")
    if ok:
        with open("tmp1.json", "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print("Saved payloads to tmp.json")
    else:
        print("Failed to parse PDF.")


    

    # workers python threads not true parallel
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    embed_pb2_grpc.add_EmbedderServicer_to_server(
        ColBERTEmbedder(embedder),
        server,
    )



    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server started on port 50051.")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()


