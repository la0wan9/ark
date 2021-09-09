package adoc

import (
	"context"
	"encoding/csv"

	"google.golang.org/grpc"

	"github.com/ahmetb/go-linq/v3"
	"github.com/la0wan9/ark/data"
	adocv1 "github.com/la0wan9/ark/pkg/adoc/v1"
	"github.com/spf13/cast"
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
		adoc := &Adoc{
			Code:   cast.ToInt64(recorder[0]),
			Parent: cast.ToInt64(recorder[1]),
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
func (s *Server) Index(ctx context.Context, req *adocv1.IndexRequest) (*adocv1.IndexResponse, error) {
	var adocMessages []*adocv1.Adoc
	linq.From(adocs).WhereT(func(a *Adoc) bool {
		ok := false
		if code := req.GetCode(); code != 0 {
			ok = true
			if code != a.Code {
				return false
			}
		}
		if parent := req.GetParent(); parent != 0 {
			ok = true
			if parent != a.Parent {
				return false
			}
		}
		if name := req.GetName(); name != "" {
			ok = true
			if name != a.Name {
				return false
			}
		}
		if !ok {
			return false
		}
		return true
	}).SelectT(func(a *Adoc) *adocv1.Adoc {
		return FromAdocToMessage(a)
	}).ToSlice(&adocMessages)
	response := &adocv1.IndexResponse{
		Adocs: adocMessages,
	}
	return response, nil
}
