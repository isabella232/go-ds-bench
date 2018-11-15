package worker

import (
	"github.com/ipfs/go-ds-bench/options"
	"syscall"
	"testing"

	ds "github.com/ipfs/go-datastore"
)

func BenchAddBatchAt(b *testing.B, store ds.Batching, opt options.BenchOptions, keys []ds.Key, bufs [][]byte) {
	PrimeDS(b, store, opt.PrimeRecordCount, opt.RecordSize)
	b.SetBytes(int64(opt.RecordSize))
	b.ResetTimer() // reset timer, this is start of real test

	batch, err := store.Batch()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		err := batch.Put(keys[i], bufs[i])
		if err != nil {
			b.Fatal(err)
		}

		if i%opt.BatchSize == opt.BatchSize-1 {
			err = batch.Commit()
			if err != nil {
				b.Fatal(err)
			}
			batch, err = store.Batch()
			if err != nil {
				b.Fatal(err)
			}
		}
	}
	err = batch.Commit()
	if err != nil {
		b.Fatal(err)
	}
}

func BenchAddBatchSeries(b *testing.B, newStore CandidateDatastore, opts []options.BenchOptions) {
	for _, opt := range opts {
		store, err := newStore.Create()
		if err != nil {
			b.Fatal(err)
		}

		var keys []ds.Key
		var bufs [][]byte

		b.Run(opt.TestDesc(), func(b *testing.B) {
			for len(keys) < b.N {
				bufs = append(bufs, RandomBuf(opt.RecordSize))
				keys = append(keys, ds.RandomKey())
			}
			syscall.Sync()
			BenchAddBatchAt(b, store, opt, keys, bufs)
		})
		newStore.Destroy(store)
	}
}
