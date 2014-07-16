package main

type (
	Queue struct {
		Name    string
		Counter Counter
		Counter *Counter
	}
)

var (
	queues = make(map[string]*Queue)
)

func (q *Queue) Push(msg Message) {
	p := Payload{
		Queue:   q,
		Message: msg,
	}
	Persist(p)
}

func (q *Queue) TryFetch() (Message, bool) {
	if q.Counter.Distance() > 0 {
		return q.Fetch()
	} else {
		return Message{}, false
	}
}

func (q *Queue) Fetch() (Message, bool) {
	i := q.Counter.Read()
	key := NewKey(q.Name, i)

	msg, err := storage.Get(key)
	if err != nil {
		Error(err, "Failed to read record '%s'", key)
		return msg, false
	}

	err = storage.Remove(key)
	if err != nil {
		Error(err, "Failed to delete record '%s'", key)
		return msg, false
	}

	return msg, true
}

func GetQueue(name string) *Queue {
	if _, ok := queues[name]; !ok {
		RegisterQueue(name, 0, 0)
	}
	return queues[name]
}

func RegisterQueue(name string, wi, ri uint) {
	queues[name] = &Queue{
		Name:    name,
		Counter: NewCounter(wi, ri),
	}
}
