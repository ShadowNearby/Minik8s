package storage

import "github.com/enriquebris/goconcurrentqueue"

// TaskQueue use TaskQueue to store etcd operations, a go function will in background to do this
var TaskQueue = goconcurrentqueue.NewFIFO()
