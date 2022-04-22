package shorten

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/urlmapping"
	pb "nilenso.com/chinnaswamy/urlmappingpb"
	"strconv"
)

type UrlMappingStore interface {
	StoreUrlMapping(ctx context.Context, mapping urlmapping.UrlMapping) error
}

type ShorteningService struct {
	pb.UnimplementedUrlShortenerServer
	store UrlMappingStore
}

func NewShorteningService(store UrlMappingStore) *ShorteningService {
	return &ShorteningService{store: store}
}

func genShortUrl() string {
	return strconv.Itoa(rand.Int())
}

func (s *ShorteningService) ShortenUrl(ctx context.Context, shortenRequest *pb.ShortenUrlRequest) (*pb.UrlMapping, error) {
	longUrl := shortenRequest.GetLongUrl()
	shortUrl := genShortUrl()
	ttl := shortenRequest.GetTTL()
	urlMapping := urlmapping.UrlMapping{ShortUrl: shortUrl, LongUrl: longUrl, TTL: ttl.AsDuration()}

	err := s.store.StoreUrlMapping(ctx, urlMapping)
	if err != nil {
		return nil, err
	}

	return urlMapping.AsProtoBuf(), nil
}

func (s *ShorteningService) Serve(ctx context.Context, done chan struct{}) {
	lis, err := net.Listen("tcp", config.ShortenListenAddress())
	if err != nil {
		log.Errorw("Could not listen on the given address",
			"listenAddress", config.ShortenListenAddress(),
			"errorMessage", err,
		)
	}
	server := grpc.NewServer()
	pb.RegisterUrlShortenerServer(server, s)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		log.Infow("Shorten server has been shut down gracefully")
	}()

	log.Infow("Starting shorten server",
		"listenAddress", config.ShortenListenAddress(),
	)
	if err := server.Serve(lis); err != nil {
		log.Errorw("Shorten server failed to serve", "errorMessage", err)
		done <- struct{}{}
		return
	}

	done <- struct{}{}
}
