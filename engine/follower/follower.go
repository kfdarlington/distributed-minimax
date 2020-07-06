package follower

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
)

type minimaxServer struct {
	pb.UnimplementedMinimaxServer
}

func (server *minimaxServer) GetExpansion(req *pb.ExpandRequest, stream pb.Minimax_GetExpansionServer) error {
	reply := &pb.ExpandReply{
		Board: req.Board,
	}
	if err := stream.Send(reply); err != nil {
		return err
	}
	return nil
}

func (server *minimaxServer) GetEvaluation(ctx context.Context, req *pb.EvaluateRequest) (*pb.EvaluateReply, error) {
	return &pb.EvaluateReply{
		Score: 500,
	}, nil
}

func (server *minimaxServer) RequestCancellation(ctx context.Context, req *pb.CancelRequest) (*pb.CancelAck, error) {
	return nil, nil
}

func NewFollower() *minimaxServer {
	return &minimaxServer{}
}
