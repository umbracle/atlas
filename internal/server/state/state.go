package state

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/umbracle/atlas/internal/proto"
	gproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	nodeBucket      = []byte("node")
	eventsBucket    = []byte("events")
	volumeBucket    = []byte("volume")
	providersBucket = []byte("providers")
)

type State struct {
	db *bolt.DB
}

func NewState(path string) (*State, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			nodeBucket,
			volumeBucket,
			providersBucket,
		}
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	state := &State{
		db: db,
	}
	return state, nil
}

func (s *State) Close() error {
	return s.db.Close()
}

func (s *State) AddNodeEvent(nodeID string, msg string) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists(eventsBucket)
	if err != nil {
		return err
	}
	nodeBkt, err := bkt.CreateBucketIfNotExists([]byte(nodeID))
	if err != nil {
		return err
	}

	id, _ := nodeBkt.NextSequence()

	event := &proto.NodeEvent{
		Message:   msg,
		Timestamp: timestamppb.New(time.Now()),
	}
	raw, err := gproto.Marshal(event)
	if err != nil {
		return err
	}
	if err := nodeBkt.Put(itob(int(id)), raw); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *State) GetNodeEvents(nodeID string) ([]*proto.NodeEvent, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt := tx.Bucket(eventsBucket)
	if bkt == nil {
		return []*proto.NodeEvent{}, nil
	}
	nodeBkt := bkt.Bucket([]byte(nodeID))
	if nodeBkt == nil {
		return []*proto.NodeEvent{}, nil
	}

	events := []*proto.NodeEvent{}
	err = nodeBkt.ForEach(func(k, v []byte) error {
		event := proto.NodeEvent{}
		if err := gproto.Unmarshal(v, &event); err != nil {
			return err
		}
		events = append(events, &event)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (s *State) UpsertNode(node *proto.Node) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists(nodeBucket)
	if err != nil {
		return err
	}
	data, err := gproto.Marshal(node)
	if err != nil {
		return err
	}
	if err := bkt.Put([]byte(node.Id), data); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *State) LoadNode(id string) (*proto.Node, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt := tx.Bucket(nodeBucket)

	raw := bkt.Get([]byte(id))
	if len(raw) == 0 {
		return nil, fmt.Errorf("not found")
	}
	var node proto.Node
	if err := gproto.Unmarshal(raw, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (s *State) ListNodes() ([]*proto.Node, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt := tx.Bucket(nodeBucket)

	nodes := []*proto.Node{}
	bkt.ForEach(func(key, val []byte) error {
		var node proto.Node
		if err := gproto.Unmarshal(val, &node); err != nil {
			return err
		}
		nodes = append(nodes, &node)
		return nil
	})

	return nodes, nil
}

func (s *State) CreateProvider(p *proto.Provider) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists(providersBucket)
	if err != nil {
		return err
	}
	data, err := gproto.Marshal(p)
	if err != nil {
		return err
	}
	if err := bkt.Put([]byte(p.Id), data); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *State) ListProviders() ([]*proto.Provider, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt := tx.Bucket(providersBucket)

	providers := []*proto.Provider{}
	bkt.ForEach(func(key, val []byte) error {
		var provider proto.Provider
		if err := gproto.Unmarshal(val, &provider); err != nil {
			return err
		}
		providers = append(providers, &provider)
		return nil
	})

	return providers, nil
}
