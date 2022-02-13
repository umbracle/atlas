package state

import (
	"github.com/boltdb/bolt"
	"github.com/umbracle/atlas/internal/proto"
	gproto "google.golang.org/protobuf/proto"
)

var (
	nodeBucket   = []byte("node")
	volumeBucket = []byte("volume")
)

type State struct {
	db *bolt.DB
}

func NewState(path string) (*State, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			nodeBucket,
			volumeBucket,
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

func (s *State) ListNodes() ([]*proto.Node, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}

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
