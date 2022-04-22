package urlmapping

import (
	"google.golang.org/protobuf/types/known/durationpb"
	pb "nilenso.com/chinnaswamy/urlmappingpb"
	"time"
)

type UrlMapping struct {
	ShortUrl string
	LongUrl  string
	TTL      time.Duration
}

func (u *UrlMapping) FromProtoBuf(urlMappingPb *pb.UrlMapping) *UrlMapping {
	u.LongUrl = urlMappingPb.GetLongUrl()
	u.ShortUrl = urlMappingPb.GetShortUrl()
	u.TTL = urlMappingPb.GetTTL().AsDuration()
	return u
}

func (u *UrlMapping) AsProtoBuf() *pb.UrlMapping {
	return &pb.UrlMapping{
		ShortUrl:  u.ShortUrl,
		LongUrl:   u.LongUrl,
		Valid:     u.TTL > 0,
		TTL:       durationpb.New(u.TTL),
		CreatedAt: nil,
	}
}
