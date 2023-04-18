package broker

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNoopBroker(t *testing.T) {
	ch := make(chan *User)

	b := NewNoopBroker()

	b.Subscribe("test", func(e Event) error {
		m := e.Message()
		u := &User{}
		_ = json.Unmarshal(m.Body, u)
		ch <- u
		return nil
	})

	go func() {
		time.Sleep(3 * time.Second)
		buf, _ := json.Marshal(&User{Name: "jack", Age: 21})
		b.Publish("test", &Message{Body: buf})

		buf, _ = json.Marshal(&User{Name: "rose", Age: 30})
		b.Publish("test", &Message{Body: buf})
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
