package cosmos_governance_bot

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetGrpcMessage(grpcUrl string, offset int64) (*v1.QueryProposalsResponse, error) {
	grpcConnection, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer grpcConnection.Close()

	client := v1.NewQueryClient(grpcConnection)
	response, err := client.Proposals(context.Background(), &v1.QueryProposalsRequest{
		Pagination: &query.PageRequest{
			Limit:  100,
			Offset: uint64(offset),
		},
	})

	if err != nil {
		log.Printf("Error getting proposals: %s", err)
	}

	return response, err
}
