package batch

import (
	"log"
	"sync"
	"time"
)

type user struct {
	ID int64
}

//add user into list
func addUser(list *[]user, ch <-chan int64, wg *sync.WaitGroup, mx *sync.Mutex) {
	defer wg.Done()

	for num := range ch {
		user := getOne(num)
		mx.Lock()
		*list = append(*list, user)
		mx.Unlock()
	}

}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	var (
		userList []user
		wg       sync.WaitGroup
		mx       sync.Mutex
	)

	//it`s unuseful: add more theads than users`
	if pool > n {
		log.Fatal("pool cannot be bigger than count of users!")
	}

	nChan := make(chan int64)

	//init goroutines
	wg.Add(int(pool))
	for i := 1; i <= int(pool); i++ {
		go addUser(&userList, nChan, &wg, &mx)
	}

	//write into channel
	for j := 0; j < int(n); j++ {
		nChan <- int64(j)
	}
	close(nChan)

	wg.Wait()

	return userList
}
