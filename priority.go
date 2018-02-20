package main

import (
	"fmt"
)

type msg_1 struct {
}

func (msg_1) msgType() {
	fmt.Println("msg type 1")
}

type msg_2 struct {
}

func (msg_2) msgType() {
	fmt.Println("msg type 2")
}

func main() {

	msg_1ch := make(chan msg_1, 49)
	msg_2ch := make(chan msg_2, 56)
	defer close(msg_1ch)
	defer close(msg_2ch)

	for i := 0; i < 49; i++ {
		msg_1ch <- msg_1{}
	}
	for i := 0; i < 56; i++ {
		msg_2ch <- msg_2{}
	}

	fmt.Println("Ordered by Channel Select Semantics")
	consumeAccordingToChannelSemantics(msg_1ch, msg_2ch)
	fmt.Println("Ordered by Nil Bias Function")
	consumeOrderedByBias(msg_1ch, msg_2ch)
	fmt.Println("Ordered by Affecting Select Probability")
	consumeOrderedByAffectingProbability(msg_1ch, msg_2ch)
}

func consumeAccordingToChannelSemantics(msg_1ch chan msg_1, msg_2ch chan msg_2) {
	// because the channels were asychronous and completely full
	// before running this routine, select's channel semantics
	// considers the messages as having arrived at the same time.
	// select therefore randomly reads one of the channels. since
	// only two case statements, probability of selection is
	// equivalent to a random coin toss.
	for msgCnt := 0; msgCnt < 21; msgCnt++ {
		select {
		case msg, ok := <-msg_1ch:
			if ok {
				msg.msgType()
			}
		case msg, ok := <-msg_2ch:
			if ok {
				msg.msgType()
			}
		}
	}
}

func consumeOrderedByBias(msg_1ch chan msg_1, msg_2ch chan msg_2) {
	// copy channels to enable their restoration
	msg_1_save := msg_1ch
	msg_2_save := msg_2ch

	// bias function encoded as a for loop
	for msgCnt := 0; msgCnt < 21; msgCnt++ {
		//	use modulus math to help implement bias function
		if msgCnt%3 == 0 {
			// favor channel 1 when processing muliples of 3
			msg_1ch = msg_1_save
			// bias channel 2
			msg_2ch = nil
		} else {
			// favor channel 2 when not a multiple of 3
			msg_2ch = msg_2_save
			// bias channel 1
			msg_1ch = nil
		}

		select {
		case msg, ok := <-msg_1ch:
			if ok {
				msg.msgType()
			}
		case msg, ok := <-msg_2ch:
			if ok {
				msg.msgType()
			}
		}
	}
}
func consumeOrderedByAffectingProbability(msg_1ch chan msg_1, msg_2ch chan msg_2) {
	// implement priorty function that selects messages from msg_2ch twice as often
	// as msg_1ch by affecting the probability distribution of the select
	// statement. produces outcome whose message type totals are equivalent to
	// nil bias function but using simpler encoding.
	for msgCnt := 0; msgCnt < 21; msgCnt++ {
		select {
		case msg, ok := <-msg_1ch:
			if ok {
				msg.msgType()
			}
		//case <-pop:
		case msg, ok := <-msg_2ch:
			if ok {
				msg.msgType()
			}
		case msg, ok := <-msg_2ch:
			if ok {
				msg.msgType()
			}
		}
	}
}
func consumeOverhead(msg_2ch <-chan msg_2, done chan<- bool) {
	defer close(done)
	// determine overhead of probability function
	for msgCnt := 0; msgCnt < 100000; msgCnt++ {
		select {
		case <-msg_2ch:
		case <-msg_2ch:
		}
	}
}
func consumeBaseLine(msg_2ch <-chan msg_2, done chan<- bool) {
	// provide baseline without probability function
	defer close(done)

	for msgCnt := 0; msgCnt < 100000; msgCnt++ {
		select {
		case <-msg_2ch:
		}
	}
}
