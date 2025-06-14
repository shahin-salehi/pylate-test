import grpc
import embed_pb2
import embed_pb2_grpc

def run():
    # Connect to gRPC server
    channel = grpc.insecure_channel('localhost:50051')
    stub = embed_pb2_grpc.EmbedderStub(channel)

    # Create request
    request = embed_pb2.EmbedRequest(text="transition plan")

    # Call RPC
    response = stub.Embed(request)

    # Display results
    for match in response.result:
        print(f"📄 {match.filename} - Page {match.page_number} - Score: {match.score:.2f}")
        if match.category:
            print(f"Category: {match.category}")
        if match.html:
            print("[TABLE CONTENT]")
            print(match.html)
        else:
            print(match.content)
        print("-" * 60)

if __name__ == "__main__":
    run()
