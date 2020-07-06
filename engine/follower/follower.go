package follower

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io/ioutil"
	"math/rand"
)

type minimaxServer struct {
	logger *logger.Logger
	pb.UnimplementedMinimaxServer
}

func (s *minimaxServer) GetExpansion(req *pb.ExpandRequest, stream pb.Minimax_GetExpansionServer) error {
	for i := 0; i < 2; i++ { // TODO fix
		reply := &pb.ExpandReply{
			Board: req.Board,
		}
		s.logger.Infof("send %d", i)
		if err := stream.Send(reply); err != nil {
			s.logger.Errorf("error when sending on stream err=%v", err)
			return err
		}
	}
	return nil
}

func (s *minimaxServer) GetEvaluation(ctx context.Context, req *pb.EvaluateRequest) (*pb.EvaluateReply, error) {
	return &pb.EvaluateReply{
		Score: rand.Float32(),
	}, nil
}

func (s *minimaxServer) RequestCancellation(ctx context.Context, req *pb.CancelRequest) (*pb.CancelAck, error) {
	return nil, nil
}

func NewFollower() *minimaxServer {
	return &minimaxServer{
		logger: logger.Init("Follower", true, false, ioutil.Discard),
	}
}
