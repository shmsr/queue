package extqueue

import (
	"strconv"
	"testing"

	"github.com/loov/queue/testsuite"
)

func Test(t *testing.T) {
	t.Helper()
	for _, desc := range Descs {
		batchSizes := testsuite.BatchSizes
		if !desc.BatchSize() {
			batchSizes = []int{0}
		}

		testSizes := testsuite.TestSizes
		if desc.Unbounded() {
			testSizes = []int{0}
		}

		t.Run(desc.Name, func(t *testing.T) {
			t.Helper()
			for _, batchSize := range batchSizes {
				for _, size := range testSizes {
					if size <= batchSize {
						continue
					}

					name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
					t.Run(name, func(t *testing.T) {
						t.Helper()
						testsuite.Tests(t, func() testsuite.Queue {
							return desc.Create(batchSize, size)
						})
					})
				}
			}
		})
	}
}

func Benchmark(b *testing.B) {
	b.Helper()
	for _, desc := range Descs {
		batchSizes := testsuite.BenchBatchSizes
		if !desc.BatchSize() {
			batchSizes = []int{0}
		}

		benchSizes := testsuite.BenchSizes
		if desc.Unbounded() {
			benchSizes = []int{0}
		}

		b.Run(desc.Name, func(b *testing.B) {
			b.Helper()
			for _, batchSize := range batchSizes {
				for _, size := range benchSizes {
					if size <= batchSize {
						continue
					}

					name := "b" + strconv.Itoa(batchSize) + "s" + strconv.Itoa(size)
					b.Run(name, func(b *testing.B) {
						b.Helper()
						testsuite.Benchmarks(b, func() testsuite.Queue {
							return desc.Create(batchSize, size)
						})
					})
				}
			}
		})
	}
}
