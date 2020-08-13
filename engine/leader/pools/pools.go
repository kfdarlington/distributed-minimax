package pools

import (
	"context"
	"errors"
	"fmt"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type follower struct {
	id int
	addr string
	client *pb.MinimaxClient
	conn *grpc.ClientConn
}

type Pool struct {
	mu sync.Mutex
	idleChan chan *follower
	followers []*follower
}

func (f *follower) GetClient() *pb.MinimaxClient {
	return f.client
}

// search for follower id in active pool - if found, move from active to idle, else raise error
func (p *Pool) MarkAsIdle(follower *follower) error {
	select {
	case p.idleChan <- follower:
		return nil
	default:
		return errors.New("idle channel was full when attempting to push follower onto it")
	}
}

func (p *Pool) Activate(ctx context.Context) *follower {
	select {
	case follower := <-p.idleChan:
		return follower
	case <-ctx.Done():
		return nil
	}
}

func (p *Pool) AddFollower(addr string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	// check if a follower in the pools already exists with address
	for _, follower := range p.followers {
		if follower.addr == addr {
			return errors.New(fmt.Sprintf("follower already exists in idle pool addr=%s", addr))
		}
	}
	// attempt to connect with follower
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(1000)*time.Millisecond)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect addr=%s err=%v", addr, err))
	} else {
		// add follower to pools
		client := pb.NewMinimaxClient(conn)
		f := &follower{
			id: len(p.followers),
			addr: addr,
			client: &client,
			conn: conn,
		}
		// increment size of idle chan
		p.followers = append(p.followers, f)
		p.idleChan <- f
	}
	return nil
}

func (p *Pool) GetFollowerAddresses() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	addresses := make([]string, len(p.followers))
	for i, follower := range p.followers {
		addresses[i] = follower.addr
	}
	return addresses
}

func (p *Pool) DestroyConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, follower := range p.followers {
		_ = follower.conn.Close()
	}
}

// assumes that followers provided to function are initially idle
func CreatePool() *Pool {
	return &Pool{
		idleChan: make(chan *follower, 1000), // TODO: figure out how to grow the idleChan -- for now, it is limited at 1000 followers
		followers: make([]*follower, 0),
	}
}
