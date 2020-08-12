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

type pool struct {
	mu sync.Mutex
	followers  []*follower // channel containing clients - number of clients in channel can be found with len(clients)
	arrived chan *follower      // a scheduler can listen to this channel to be notified of a new client being added
}

type Pools struct {
	mu sync.Mutex
	idle *pool
	active *pool
}

func (f *follower) GetClient() *pb.MinimaxClient {
	return f.client
}

func (p *pool) pop() *follower {
	p.mu.Lock()
	defer p.mu.Unlock()
	count := len(p.followers)
	if count < 1 {
		// signifies empty pool
		return nil
	}
	// choose client from pool that has been in the pool the longest
	f := p.followers[0]
	// these next 3 lines are a temporary workaround because the original method of
	// p.followers = p.followers[1:] decreases the original capacity by 1 and breaks things
	newFollowers := make([]*follower, len(p.followers) - 1, cap(p.followers))
	copy(newFollowers, p.followers[1:])
	p.followers = newFollowers
	return f
}

func (p *pool) push(f *follower) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	count := len(p.followers)
	if count == cap(p.followers) {
		return errors.New("entered follower into full pool")
	}
	select {
	case p.arrived <- f: // signal to a listener that a client was added
	default: // nobody was listening - do not block
		p.followers = p.followers[:count + 1]
		p.followers[count] = f
	}
	return nil
}

func (p *pool) remove(f *follower) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	idx := -1
	for i, follower := range p.followers {
		if follower.id == f.id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("follower was not in pool")
	}
	p.followers = append(p.followers[:idx], p.followers[idx + 1:]...)
	return nil
}

func createPool() *pool {
	return &pool{
		followers: make([]*follower, 0),
		arrived: make(chan *follower),
	}
}

// search for follower id in active pool - if found, move from active to idle, else raise error
func (p *Pools) MarkIdle(follower *follower) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.active.remove(follower); if err != nil {
		return err
	}
	err = p.idle.push(follower); if err != nil {
		return err
	}
	return nil
}

// search for follower id in idle pool - if found, move from idle to active, else raise error
func (p *Pools) MarkActive(follower *follower) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	err := p.idle.remove(follower); if err != nil {
		return err
	}
	err = p.active.push(follower); if err != nil {
		return err
	}
	return nil
}

func (p *Pools) Activate(ctx context.Context) (*follower, error) {
	// try to retrieve a follower right away
	p.mu.Lock()
	follower := p.idle.pop()
	p.mu.Unlock()

	// if no idle followers, listen for the next available follower until context expires
	if follower == nil {
		select {
		case follower = <-p.idle.arrived:
		case <-ctx.Done():
			return nil, nil
		}
	}
	if err := p.active.push(follower); err != nil {
		return follower, err
	}
	return follower, nil
}

func (p *Pools) AddFollower(addr string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	// check if a follower in the pools already exists with address
	for _, follower := range p.idle.followers {
		if follower.addr == addr {
			return errors.New(fmt.Sprintf("follower already exists in idle pool addr=%s", addr))
		}
	}
	for _, follower := range p.active.followers {
		if follower.addr == addr {
			return errors.New(fmt.Sprintf("follower already exists in active pool addr=%s", addr))
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
			id: len(p.idle.followers) + len(p.active.followers) + 1,
			addr: addr,
			client: &client,
			conn: conn,
		}
		// increase capacity of pools by 1
		p.idle.mu.Lock()
		newIdleFollowers := make([]*follower, len(p.idle.followers), cap(p.idle.followers) + 1)
		copy(newIdleFollowers, p.idle.followers)
		p.idle.followers = newIdleFollowers
		p.idle.mu.Unlock()
		p.active.mu.Lock()
		newActiveFollowers := make([]*follower, len(p.active.followers), cap(p.active.followers) + 1)
		copy(newActiveFollowers, p.active.followers)
		p.active.followers = newActiveFollowers
		p.active.mu.Unlock()
		_ = p.idle.push(f)
	}
	return nil
}

func (p *Pools) DestroyConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idle.mu.Lock()
	for _, follower := range p.idle.followers {
		_ = follower.conn.Close()
	}
	p.idle.mu.Unlock()
	p.active.mu.Lock()
	for _, follower := range p.active.followers {
		_ = follower.conn.Close()
	}
	p.active.mu.Unlock()
}

// assumes that followers provided to function are initially idle
func CreatePools() *Pools {
	// create idle pool
	idlePool := createPool()

	// create active pool
	activePool := createPool() // error would have been raised above

	return &Pools{
		idle: idlePool,
		active: activePool,
	}
}
