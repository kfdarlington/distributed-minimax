package pools

import (
	"errors"
	"github.com/kristian-d/distributed-battlesnake/minimax/pb"
	"sync"
)

type follower struct {
	id int
	client *pb.MinimaxClient
}

type pool struct {
	mu sync.Mutex
	followers  []*follower // channel containing clients - number of clients in channel can be found with len(clients)
	arrived chan bool      // a scheduler can listen to this channel to be notified of a new client being added
}

type Pools struct {
	mu sync.Mutex
	idle *pool
	active *pool
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
	follower := p.followers[0]
	p.followers = p.followers[1:]
	return follower
}

func (p *pool) push(f *follower) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	count := len(p.followers)
	if count == cap(p.followers) {
		return errors.New("entered follower into full pool")
	}
	p.followers = p.followers[:count + 1]
	p.followers[count] = f
	select {
	case p.arrived <- true: // signal to a listener that a client was added
	default: // nobody was listening - do not block
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

func createPool(cap int) (*pool, error) {
	if cap < 1 {
		return nil, errors.New("cannot create pool with capacity less than 1")
	}
	return &pool{
		followers: make([]*follower, 0, cap),
		arrived: make(chan bool),
	}, nil
}

// search for follower id in active pool - if found, move from active to idle, else raise error
func (p *Pools) markIdle(follower *follower) error {
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
func (p *Pools) markActive(follower *follower) error {
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

func (p *Pools) activate() (*follower, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	follower := p.idle.pop(); if follower == nil {
		return nil, nil
	}
	err := p.active.push(follower); if err != nil {
		return follower, err
	}
	return follower, nil
}

// assumes that followers provided to function are initially idle
func CreateFollowerPools(clients []*pb.MinimaxClient) (*Pools, error) {
	// create followers from clients
	followers := make([]*follower, len(clients))
	for i, client := range clients {
		followers[i] = &follower{
			id: i,
			client: client,
		}
	}

	// create idle pool
	idlePool, err := createPool(len(followers))
	if err != nil {
		return nil, err
	}
	for _, follower := range followers {
		_ = idlePool.push(follower) // error would have been raised above
	}

	// create active pool
	activePool, _ := createPool(len(followers)) // error would have been raised above

	return &Pools{
		idle: idlePool,
		active: activePool,
	}, nil
}
