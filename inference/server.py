
import grpc 
from concurrent import futures
import onnxruntime as ort 
from transformers import AutoTokenizer
import numpy as np
import embed_pb2, embed_pb2_grpc 


class ColBERTEmbedder(embed_pb2_grpc.EmbedderServicer):
    def __init__(self, model_path: str, tokenizer_name: str):
        self.session = ort.InferenceSession(model_path, providers=["CPUExecutionProvider"])
        self.tokenizer = AutoTokenizer.from_pretrained(tokenizer_name)

    def Embed(self, request, context):
        # Tokenize input
        tokens = self.tokenizer(
            request.text,
            return_tensors="np",
            padding="max_length",
            truncation=True,
            max_length=128,
        )

        # ONNX expects inputs as dict of numpy arrays
        onnx_inputs = {
            "input_ids": tokens["input_ids"],
            "attention_mask": tokens["attention_mask"]
        }
        outputs = self.session.run(None, onnx_inputs)
        
        # Assume output[0] is the embedding
        embedding = outputs[0].mean(axis=1).flatten().tolist()  # e.g., [batch, dim] → pooled 
        return embed_pb2.EmbedResponse(embedding=embedding)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    embed_pb2_grpc.add_EmbedderServicer_to_server(
        ColBERTEmbedder("model/model.onnx", "bert-base-uncased"),
        server,
    )
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server started on port 50051.")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()


