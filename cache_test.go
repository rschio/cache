package cache

import "testing"

var ss = []string{"hello", "world", "foo", "bar", "one", "day",
	"light", "time", "month", "car"}

func createAndPop() *Cache {
	c := New(10)
	for i, word := range ss {
		c.Insert(word, i)
	}
	return c
}

func TestNew(t *testing.T) {
	c := New(0)
	c.Insert(ss[0], 10)
	c.Insert(ss[1], 11)
	_, ok := c.Get(ss[0])
	if ok {
		t.Error("got value that should be replaced")
	}
	el, ok := c.Get(ss[1])
	if !ok {
		t.Error("failed to get value")
	}
	if el.(int) != 11 {
		t.Errorf("got wrong value")
	}
}

func TestSetMax(t *testing.T) {
	c := createAndPop()
	c.SetMax(2)
	items := make([]item, 0)
	fn := fnRange(&items)
	c.Range(fn)

	if len(items) != 2 {
		t.Error("wrong number of members")
	}
}

func TestGet(t *testing.T) {
	c := createAndPop()
	for i, word := range ss {
		el, ok := c.Get(word)
		if !ok {
			t.Errorf("failed to get word: %s", word)
		}
		if el.(int) != i {
			t.Error("got wrong value")
		}
	}
}

func TestDelete(t *testing.T) {
	c := createAndPop()
	_, ok := c.Get(ss[0])
	if !ok {
		t.Error("get failed to get value")
	}
	c.Delete(ss[0])
	_, ok = c.Get(ss[0])
	if ok {
		t.Error("got value that should be deleted")
	}
}

type item struct {
	key string
	val interface{}
}

func fnRange(items *[]item) func(string, interface{}) {
	return func(s string, el interface{}) {
		i := item{key: s, val: el}
		*items = append(*items, i)
	}
}

func TestRange(t *testing.T) {
	c := createAndPop()
	items := make([]item, 0)
	fn := fnRange(&items)
	c.Range(fn)
	if len(items) < 1 {
		t.Error("got len < 1")
	}
	for _, it := range items {
		if ss[it.val.(int)] != it.key {
			t.Error("got wrong item value in range")
		}
	}
}
