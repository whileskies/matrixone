package spillEngine

import (
	"matrixone/pkg/encoding"
	"matrixone/pkg/vm/engine"
	"matrixone/pkg/vm/engine/spillEngine/kv"
	"matrixone/pkg/vm/engine/spillEngine/meta"
	"matrixone/pkg/vm/metadata"
	"os"
	"path"
	"runtime"
)

const (
	MetaKey = "meta"
)

func New(path string, db engine.DB) (*spillEngine, error) {
	cdb, err := kv.New(path)
	if err != nil {
		return nil, err
	}
	return &spillEngine{db: db, cdb: cdb, path: path}, nil
}

func (e *spillEngine) Close() error {
	return e.db.Close()
}

func (e *spillEngine) Del(k []byte) error {
	return e.db.Del(k)
}

func (e *spillEngine) Set(k, v []byte) error {
	return e.db.Set(k, v)
}

func (e *spillEngine) Get(k []byte) ([]byte, error) {
	return e.db.Get(k)
}

func (e *spillEngine) NewBatch() (engine.Batch, error) {
	return e.db.NewBatch()
}

func (e *spillEngine) NewIterator(prefix []byte) (engine.Iterator, error) {
	return e.db.NewIterator(prefix)
}

func (e *spillEngine) Node(_ string) *engine.NodeInfo {
	return &engine.NodeInfo{
		Mcpu: runtime.NumCPU(),
	}
}

func (e *spillEngine) Delete(name string) error {
	return nil
}

func (e *spillEngine) Create(name string) error {
	return nil
}

func (e *spillEngine) Databases() []string {
	return nil
}

func (e *spillEngine) Database(name string) (engine.Database, error) {
	return &database{e.path, e.cdb, e.db}, nil
}

func (e *database) Delete(name string) error {
	return os.RemoveAll(path.Join(e.path, name))
}

func (e *database) Create(name string, defs []engine.TableDef, _ *engine.PartitionBy, _ *engine.DistributionBy) error {
	var attrs []metadata.Attribute

	{
		for _, def := range defs {
			v, ok := def.(*engine.AttributeDef)
			if ok {
				attrs = append(attrs, v.Attr)
			}
		}
	}
	data, err := encoding.Encode(meta.Metadata{Name: name, Attrs: attrs})
	if err != nil {
		return err
	}
	dir := path.Join(e.path, name)
	if _, err := os.Stat(dir); os.IsExist(err) {
		return os.ErrExist
	}
	if err := os.Mkdir(dir, os.FileMode(0775)); err != nil {
		return err
	}
	if err := e.cdb.Set(path.Join(name, MetaKey), data); err != nil {
		os.RemoveAll(path.Join(e.path, name))
		return err
	}
	return nil
}

func (e *database) Relations() []string {
	return nil
}

func (e *database) Relation(name string) (engine.Relation, error) {
	var md meta.Metadata

	data, err := e.cdb.GetCopy(path.Join(name, MetaKey))
	if err != nil {
		return nil, err
	}
	if err := encoding.Decode(data, &md); err != nil {
		return nil, err
	}
	mp := make(map[string]metadata.Attribute)
	{
		for _, attr := range md.Attrs {
			mp[attr.Name] = attr
		}
	}
	return &relation{
		md: md,
		mp: mp,
		id: name,
		db: e.cdb,
	}, nil
}
