// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

import (
	dbsched "matrixone/pkg/vm/engine/aoe/storage/db/sched"
	"matrixone/pkg/vm/engine/aoe/storage/events/memdata"
	"matrixone/pkg/vm/engine/aoe/storage/layout/table/v1/iface"
	"matrixone/pkg/vm/engine/aoe/storage/metadata/v2"
	"matrixone/pkg/vm/engine/aoe/storage/sched"
)

type createBlkEvent struct {
	dbsched.BaseEvent

	// TableID is Table's id, aoe is generated when the table is created
	TableID uint64

	// TableData is Table's metadata in memory
	TableData iface.ITableData

	// Block created by NewCreateBlkEvent
	Block iface.IBlock

	PrevMeta *metadata.Block
}

// NewCreateBlkEvent creates a logical Block event
func NewCreateBlkEvent(ctx *dbsched.Context, tid uint64, prevBlock *metadata.Block, tableData iface.ITableData) *createBlkEvent {
	e := &createBlkEvent{
		TableData: tableData,
		TableID:   tid,
		PrevMeta:  prevBlock,
	}
	e.BaseEvent = dbsched.BaseEvent{
		Ctx:       ctx,
		BaseEvent: *sched.NewBaseEvent(e, sched.MetaCreateBlkTask, ctx.DoneCB, ctx.Waitable),
	}
	return e
}

// Return the block just created
func (e *createBlkEvent) GetBlock() *metadata.Block {
	if e.Err != nil {
		return nil
	}
	return e.Result.(*metadata.Block)
}

func (e *createBlkEvent) Execute() error {
	table := e.Ctx.Opts.Meta.Catalog.SimpleGetTable(e.TableID)
	if table == nil {
		return metadata.TableNotFoundErr
	}

	blk := table.SimpleGetOrCreateNextBlock(e.PrevMeta)
	e.Result = blk

	if e.TableData != nil {
		ctx := &memdata.Context{Opts: e.Ctx.Opts, Waitable: true}
		event := memdata.NewCreateSegBlkEvent(ctx, blk, e.TableData)
		if err := e.Ctx.Opts.Scheduler.Schedule(event); err != nil {
			return err
		}
		if err := event.WaitDone(); err != nil {
			return err
		}
		e.Block = event.Block
	}
	return nil
}