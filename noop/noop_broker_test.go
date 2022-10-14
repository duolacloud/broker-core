package noop

import (
	"testing"
	"time"

	"github.com/duolacloud/broker-core"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNoopBroker(t *testing.T) {
	ch := make(chan *User)

	b := NewBroker()
	_ = b.Connect()

	b.Subscribe("test", func(e broker.Event) error {
		ch <- e.Message().(*User)
		return nil
	}, broker.ResultType(&User{}))

	go func() {
		time.Sleep(3 * time.Second)
		b.Publish("test", &User{Name: "jack", Age: 21})
		b.Publish("test", &User{Name: "rose", Age: 30})
	}()

	u1 := <-ch
	t.Logf("%+v", u1)
	assert.Equal(t, "jack", u1.Name)
	assert.Equal(t, 21, u1.Age)

	u2 := <-ch
	t.Logf("%+v", u2)
	assert.Equal(t, "rose", u2.Name)
	assert.Equal(t, 30, u2.Age)
}
