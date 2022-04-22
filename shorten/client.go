package shorten

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/urlmapping"
	pb "nilenso.com/chinnaswamy/urlmappingpb"
	"time"
)

type Client struct {
	rpcClient pb.UrlShortenerClient
	conn      *grpc.ClientConn
	addr      string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorw("Could not connect to gRPC server",
			"errorMessage", err,
		)
		return err
	}

	c.conn = conn
	c.rpcClient = pb.NewUrlShortenerClient(conn)
	return nil
}

func (c *Client) ShortenUrl(longUrl string, ttl time.Duration) (*urlmapping.UrlMapping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if c.rpcClient == nil {
		return nil, errors.New("gRPC connection client not initialized")
	}

	urlMappingPb, err := c.rpcClient.ShortenUrl(ctx, &pb.ShortenUrlRequest{
		LongUrl: longUrl,
		TTL:     durationpb.New(ttl),
	})

	if err != nil {
		return nil, err
	}

	urlMapping := &urlmapping.UrlMapping{}
	return urlMapping.FromProtoBuf(urlMappingPb), nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Errorw("gRPC client: could not close connection gracefully",
			"errorMessage", err,
		)
	}
}
