package follower

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/expander"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io/ioutil"
	"math/rand"
)

type minimaxServer struct {
	logger *logger.Logger
	pb.UnimplementedMinimaxServer
}

func (s *minimaxServer) GetExpansion(req *pb.ExpandRequest, stream pb.Minimax_GetExpansionServer) error {
	resultChan := make(chan *pb.Board)
	go expander.Expand(stream.Context(), req.GetBoard(), req.GetIsMaximizerTurn(), resultChan)
	for {
		select {
		case board, ok := <-resultChan:
			if !ok {
				s.logger.Info("board expansion finished on time")
				return nil
			}
			reply := &pb.ExpandReply{
				Board: board,
			}
			s.logger.Info("sending board on stream")
			if err := stream.Send(reply); err != nil {
				s.logger.Errorf("error when sending board on stream err=%v", err)
				return err
			}
		case <-stream.Context().Done():
			s.logger.Info("stream context expired during board expansion")
			return nil
		}
	}
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
