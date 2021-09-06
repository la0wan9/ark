package adoc

import (
	"context"
	"encoding/csv"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/la0wan9/ark/data"
	adocv1 "github.com/la0wan9/ark/pkg/adoc/v1"
)

var adocs []*Adoc

func init() {
	file, err := data.FS.Open("adoc/adoc.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	recorders, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}
	for _, recorder := range recorders {
		if len(recorder) != 3 {
			panic("invalid data")
		}
		code, err := strconv.ParseInt(recorder[0], 10, 64)
		if err != nil {
			panic(err)
		}
		parent, err := strconv.ParseInt(recorder[1], 10, 64)
		if err != nil {
			panic(err)
		}
		adoc := &Adoc{
			Code:   code,
			Parent: parent,
			Name:   recorder[2],
		}
		adocs = append(adocs, adoc)
	}
}

// Server implements AdocServiceServer
type Server struct {
	adocv1.UnimplementedAdocServiceServer
}

// Register registers the server service on the given gRPC server.
func (s *Server) Register(server *grpc.Server) {
	adocv1.RegisterAdocServiceServer(server, s)
}

// Index returns *adocv1.Adocs and error
func (s *Server) Index(context.Context, *emptypb.Empty) (*adocv1.Adocs, error) {
	adocResponses := make([]*adocv1.Adoc, len(adocs))
	for i, adoc := range adocs {
		adocResponses[i] = FromAdocToMessage(adoc)
	}
	response := &adocv1.Adocs{
		Adocs: adocResponses,
	}
	return response, nil
}
