package testsuite

import (
	"strconv"
	"sync"
	"testing"
)

func benchCommon(b *testing.B, caps Capability, ctor func() Queue) {
	b.Run("Create", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctor()
		}
	})
}

func benchSPSC(b *testing.B, caps Capability, ctor func() Queue) {
	b.Run("Single", func(b *testing.B) {
		q := ctor().(SPSC)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v Value
			q.Send(v)
			q.Recv(&v)
		}
	})

	b.Run("Uncontended/x100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			q := ctor().(SPSC)
			for pb.Next() {
				var v Value
				for i := 0; i < 100; i++ {
					q.Send(v)
					q.Recv(&v)
				}
			}
		})
	})

	for _, work := range BenchWork {
		suffix := ""
		if work > 0 {
			suffix = "Work" + strconv.Itoa(work)
		}
		b.Run("ProducerConsumer"+suffix, func(b *testing.B) {
			q := ctor().(SPSC)
			b.ResetTimer()
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				for i := 0; i < b.N; i++ {
					var v Value
					q.Send(v)
					LocalWork(work)
				}
				wg.Done()
			}()
			go func() {
				for i := 0; i < b.N; i++ {
					var v Value
					q.Recv(&v)
					LocalWork(work)
				}
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func benchMPSC(b *testing.B, caps Capability, ctor func() Queue) {
	for _, work := range BenchWork {
		suffix := ""
		if work > 0 {
			suffix = "Work" + strconv.Itoa(work)
		}
		b.Run("ProducerConsumer"+suffix+"/x100", func(b *testing.B) {
			q := ctor().(MPSC)
			b.ResetTimer()
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for i := 0; i < 100; i++ {
							q.Send(0)
							LocalWork(work)
						}
					}
				})
				wg.Done()
			}()

			go func() {
				for i := 0; i < b.N; i++ {
					for i := 0; i < 100; i++ {
						var v Value
						q.Recv(&v)
						LocalWork(work)
					}
				}
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func benchSPMC(b *testing.B, caps Capability, ctor func() Queue) {
	for _, work := range BenchWork {
		suffix := ""
		if work > 0 {
			suffix = "Work" + strconv.Itoa(work)
		}
		b.Run("ProducerConsumer"+suffix+"/x100", func(b *testing.B) {
			q := ctor().(SPMC)
			b.ResetTimer()
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				for i := 0; i < b.N; i++ {
					for i := 0; i < 100; i++ {
						q.Send(0)
						LocalWork(work)
					}
				}
				wg.Done()
			}()

			go func() {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for i := 0; i < 100; i++ {
							var v Value
							q.Recv(&v)
						}
					}
				})
				wg.Done()
			}()

			wg.Wait()
		})
	}
}

func benchMPMC(b *testing.B, caps Capability, ctor func() Queue) {
	b.Run("Contended/x100", func(b *testing.B) {
		q := ctor().(MPMC)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var v Value
				for i := 0; i < 100; i++ {
					q.Send(v)
					q.Recv(&v)
				}
			}
		})
	})

	for _, work := range BenchWork {
		suffix := ""
		if work > 0 {
			suffix = "Work" + strconv.Itoa(work)
		}
		b.Run("ProducerConsumer"+suffix+"/x100", func(b *testing.B) {
			q := ctor().(MPMC)
			b.ResetTimer()

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for i := 0; i < 100; i++ {
							q.Send(0)
							LocalWork(work)
						}
					}
				})
				wg.Done()
			}()

			go func() {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for i := 0; i < 100; i++ {
							var v Value
							q.Recv(&v)
							LocalWork(work)
						}
					}
				})
				wg.Done()
			}()
			wg.Wait()
		})
	}
}

func benchNonblockSPSC(b *testing.B, caps Capability, ctor func() NonblockingSPSC) {
	b.Run("Single", func(b *testing.B) {
		q := ctor().(NonblockingSPSC)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var v Value
			q.TrySend(v)
			q.TryRecv(&v)
		}
	})

	b.Run("Uncontended/x100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			q := ctor().(NonblockingSPSC)
			for pb.Next() {
				var v Value
				for i := 0; i < 100; i++ {
					q.TrySend(v)
					q.TryRecv(&v)
				}
			}
		})
	})
}

func benchNonblockMPSC(b *testing.B, caps Capability, ctor func() NonblockingMPSC) { b.Skip("todo") }
func benchNonblockSPMC(b *testing.B, caps Capability, ctor func() NonblockingSPMC) { b.Skip("todo") }

func benchNonblockMPMC(b *testing.B, caps Capability, ctor func() NonblockingMPMC) {
	b.Run("Contended/x100", func(b *testing.B) {
		q := ctor().(NonblockingMPMC)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var v Value
				for i := 0; i < 100; i++ {
					q.TrySend(v)
					q.TryRecv(&v)
				}
			}
		})
	})
}
