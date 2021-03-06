package pool

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"gopkg.in/alexcesaro/statsd.v2"
)

// FreeMac is the Free Mac address constant
const FreeMac = "00:00:00:00:00:00"

// FakeMac is the Fake Mac address constant
const FakeMac = "ff:ff:ff:ff:ff:ff"

// Random ip constant
const Random = 1

// OldestReleased ip constant
const OldestReleased = 2

// Backend interface
type Backend interface {
	NewDHCPPool(ctx context.Context, capacity uint64, algorithm int, StatsdClient *statsd.Client)
	ReserveIPIndex(index uint64, mac string) (string, error)
	IsFreeIPAtIndex(index uint64) bool
	GetMACIndex(index uint64) (uint64, string, error)
	GetFreeIPIndex(mac string) (uint64, string, error)
	IndexInPool(index uint64) bool
	FreeIPsRemaining() uint64
	FreeIPIndex(index uint64) error
	Capacity() uint64
	GetDHCPPool() DHCPPool
	GetIssues(macs []string) ([]string, map[uint64]string)
	Listen() bool
}

// Creater function
type Creater func(context.Context, uint64, string, int, *statsd.Client, *sql.DB) (Backend, error)

var poolLookup = map[string]Creater{
	"memory": NewMemoryPool,
	"mysql":  NewMysqlPool,
}

// Create function
func Create(ctx context.Context, poolType string, capacity uint64, name string, algorithm int, StatsdClient *statsd.Client, sql *sql.DB) (Backend, error) {
	if creater, found := poolLookup[poolType]; found {
		return creater(ctx, capacity, name, algorithm, StatsdClient, sql)
	}

	return nil, fmt.Errorf("Pool of %s not found", poolType)
}

// DHCPPool struct
type DHCPPool struct {
	lock      *sync.RWMutex
	free      map[uint64]bool
	mac       map[uint64]string
	capacity  uint64
	released  map[uint64]int64
	algorithm int
	ctx       context.Context
	statsd    *statsd.Client
}

// Timing struct
type Timing struct {
	timing statsd.Timing
}

// NewTiming struct
func (dp *DHCPPool) NewTiming() *Timing {
	if dp.statsd == nil {
		return nil
	}

	return &Timing{timing: dp.statsd.NewTiming()}
}

// Send function to add pfdhcp prefix
func (t *Timing) Send(name string) {
	if t == nil {
		return
	}

	t.timing.Send("pfdhcp." + name)
}

// Track timing for each function
func (dp *DHCPPool) timeTrack(t *Timing, name string) {
	t.Send("pfdhcp." + name)
}
