package koorosh

import (
	"context"
	"shahin/webserver/internal/types"
	pb "shahin/webserver/proto/embed"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct{
	conn *grpc.ClientConn
	client pb.EmbedderClient
}
	
func New(address string)(*Client, error){
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil{
		return nil, err
	}
	
	return &Client{
		conn: conn,
		client: pb.NewEmbedderClient(conn),
	}, nil
}
func (c *Client) Search(query string, category string) ([]types.Match, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	resp, err := c.client.Embed(ctx, &pb.EmbedRequest{Text: query, Category: category})
	if err != nil {
		return nil, err
	}

	results := make([]types.Match, len(resp.Result))
	for i, m := range resp.Result {
		results[i] = types.Match{
			Filename: m.Filename,
			PageNumber: m.PageNumber,
			Title: m.Title,
			Category: m.Category,
			Content: m.Content,
			HTML: m.Html,
			Score: m.Score,
			Meta: m.Meta,
		}
	}

	return results, nil
}

func (c *Client) Close(){
	c.conn.Close()
}
