// Code generated by rocksdb-bindings generator. DO NOT EDIT.
package legacy

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/db"
	"github.com/stackrox/rox/pkg/rocksdb"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
)

var (
	bucket = []byte("signature_integrations")
)

type Store interface {
	UpsertMany(ctx context.Context, objs []*storage.SignatureIntegration) error
	Walk(ctx context.Context, fn func(obj *storage.SignatureIntegration) error) error
}

type storeImpl struct {
	crud db.Crud
}

func alloc() proto.Message {
	return &storage.SignatureIntegration{}
}

func keyFunc(msg proto.Message) []byte {
	return []byte(msg.(*storage.SignatureIntegration).GetId())
}
func uniqKeyFunc(msg proto.Message) []byte {
	return []byte(msg.(*storage.SignatureIntegration).GetName())
}

// New returns a new Store instance using the provided rocksdb instance.
func New(db *rocksdb.RocksDB) (Store, error) {
	baseCRUD := generic.NewUniqueKeyCRUD(db, bucket, keyFunc, alloc, uniqKeyFunc, false)
	return  &storeImpl{crud: baseCRUD}, nil
}

// UpsertMany batches objects into the DB
func (b *storeImpl) UpsertMany(_ context.Context, objs []*storage.SignatureIntegration) error {
	msgs := make([]proto.Message, 0, len(objs))
	for _, o := range objs {
		msgs = append(msgs, o)
	}

	return b.crud.UpsertMany(msgs)
}

// Walk iterates over all of the objects in the store and applies the closure
func (b *storeImpl) Walk(_ context.Context, fn func(obj *storage.SignatureIntegration) error) error {
	return b.crud.Walk(func(msg proto.Message) error {
		return fn(msg.(*storage.SignatureIntegration))
	})
}
