package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEventSeq(t *testing.T) {
	m := NewManger(128, 16)
	m.RegisterA("ceshi", "asdf1", func(event *Event) {
		fmt.Println("Num:", event.Data)
	})
	go func() {
		for i := 0; i < 1000; i++ {
			m.CallA("ceshi", fmt.Sprintf("a%d", i))
		}
	}()
	for i := 0; i < 1000; i++ {
		m.CallA("ceshi", i)
	}
}

func TestEventWeired(t *testing.T) {
	m := NewManger(128, 2)
	m.RegisterA("playlist.update", "asdf1", func(event *Event) {
		fmt.Printf("%d %p, outdated: %t\n", event.Data, event, event.Outdated)
	})
	for i := 0; i < 2; i++ {
		fmt.Println("asdfsafasfasfasfasfasf")
		m.CallA("playlist.update", i)
		fmt.Println("asdfsafasfasfasfasfasf")
	}
	time.Sleep(1 * time.Second)
}
