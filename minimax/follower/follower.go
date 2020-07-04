package follower

import (
	"context"
	"github.com/kristian-d/distributed-battlesnake/minimax/pb"
)

type minimaxServer struct {
	pb.UnimplementedMinimaxServer
}

func (server *minimaxServer) GetExpansion(req *pb.ExpandRequest, stream pb.Minimax_GetExpansionServer) error {
	return nil
}

func (server *minimaxServer) GetEvaluation(ctx context.Context, req *pb.EvaluateRequest) (*pb.EvaluateReply, error) {
	return nil, nil
}

func (server *minimaxServer) RequestCancellation(ctx context.Context, req *pb.CancelRequest) (*pb.CancelAck, error) {
	return nil, nil
}

func NewFollower() *minimaxServer {
	return &minimaxServer{}
}
