package leader

import (
	"github.com/kristian-d/distributed-battlesnake/minimax/leader/pools"
	"github.com/kristian-d/distributed-battlesnake/minimax/pb"
)

type Leader struct {
	pools *pools.Pools
}

func (leader *Leader) Start() error {
	return nil
}

func NewLeader(clients []*pb.MinimaxClient) (*Leader, error) {
	p, err := pools.CreateFollowerPools(clients)
	if err != nil {
		return nil, err
	}
	return &Leader{
		pools: p,
	}, nil
}
