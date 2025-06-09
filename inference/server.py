import logging
import sys
import grpc 
import embed_pb2, embed_pb2_grpc 
import json
import numpy as np

from concurrent import futures
from internal.colbert.embedder import Embedder
from internal.db.database import Database
from internal.parser.parse import Parse


logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
        )

logger = logging.getLogger(__name__)

def make_json_safe(obj):
    if isinstance(obj, np.ndarray):
        return obj.tolist()
    elif isinstance(obj, list):
        return [make_json_safe(x) for x in obj]
    elif isinstance(obj, dict):
        return {k: make_json_safe(v) for k, v in obj.items()}
    return obj

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

    """
    # init parser
    parse = Parse(embedder)

    
    ## parse test
    data, ok = parse.pdf("internal/parser/docs/sector_policy_fossil_fuel.pdf", "test", "test")
    if ok:
        with open("tmp1.json", "w", encoding="utf-8") as f:
            json.dump(make_json_safe(data), f, indent=2, ensure_ascii=False)

            logger.info("pdf parsed succesfully.")
    else:
        logger.error("Failed to parse PDF.")
        sys.exit(1)

    ## insert test

    resp, ok = db.insert_pdf(data)
    if not ok:
        logger.error("I just shitted my pants, ong")
        sys.exit(1)
    else:
        logger.info(f"pdf inserted id: {resp} ")
    """
    ## search test
    l = db.search(query="transition plan", embedder=embedder)
    for row in l:
        print(row)
        
        print("\n\n")



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


