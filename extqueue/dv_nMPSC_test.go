package extqueue_test

import (
	"sync"
	"testing"

	"github.com/loov/queue/extqueue"
)

func BenchmarkPingPongSPSCns(b *testing.B) {
	q1, q2 := extqueue.NewSPSCnsDV(), extqueue.NewSPSCnsDV()
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < b.N; i++ {
			var v extqueue.Value
			q1.Send(v)
			q2.Recv(&v)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			var v extqueue.Value
			q1.Recv(&v)
			q2.Send(v)
		}
		wg.Done()
	}()
	wg.Wait()
}

func BenchmarkPingPongMPSCns(b *testing.B) {
	q1, q2 := extqueue.NewMPSCnsDV(), extqueue.NewMPSCnsDV()
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < b.N; i++ {
			var v extqueue.Value
			q1.Send(v)
			q2.Recv(&v)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			var v extqueue.Value
			q1.Recv(&v)
			q2.Send(v)
		}
		wg.Done()
	}()
	wg.Wait()
}

func BenchmarkPingPongMPSCnsi(b *testing.B) {
	q1, q2 := extqueue.NewMPSCnsiDV(), extqueue.NewMPSCnsiDV()
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		node := &extqueue.Node{}
		for i := 0; i < b.N; i++ {
			node.Value = extqueue.Value(i)
			q1.SendNode(node)
			node, _ = q2.RecvNode()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			node, _ := q1.RecvNode()
			q2.SendNode(node)
		}
		wg.Done()
	}()
	wg.Wait()
}
