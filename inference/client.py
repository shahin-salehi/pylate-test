import grpc 
import embed_pb2
import embed_pb2_grpc

def run():
    channel = grpc.insecure_channel("localhost:50051")
    stub = embed_pb2_grpc.EmbedderStub(channel)
    request = embed_pb2.EmbedRequest(text="hello world")
    response = stub.Embed(request)

    # Show the embedding result
    print(f"Received embedding of length {len(response.embedding)}")
    print(response.embedding[:10], "...")  # Show first 10 values

if __name__ == "__main__":
    run()
