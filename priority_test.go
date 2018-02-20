package main

import "testing"

func BenchmarkBaseline(b *testing.B) {
	msg_2ch := make(chan msg_2)
	done := startBaseLine(msg_2ch)
	type2MessageGen(msg_2ch)
	<-done
	close(msg_2ch)
}

func BenchmarkTestOverhead(b *testing.B) {
	msg_2ch := make(chan msg_2)
	done := startOverhead(msg_2ch)
	type2MessageGen(msg_2ch)
	<-done
	close(msg_2ch)
}

func type2MessageGen(msg_2ch chan msg_2) {
	for msgCnt := 0; msgCnt < 100000; msgCnt++ {
		msg_2ch <- msg_2{}
	}
}

func startBaseLine(msg_2ch <-chan msg_2) (done <-chan bool) {
	d := make(chan bool)
	done = d
	go consumeBaseLine(msg_2ch, d)
	return
}
func startOverhead(msg_2ch <-chan msg_2) (done <-chan bool) {
	d := make(chan bool)
	done = d
	go consumeOverhead(msg_2ch, d)
	return
}
