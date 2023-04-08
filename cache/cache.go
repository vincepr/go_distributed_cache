package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache{
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, val []byte, ttl time.Duration) error{
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = val
	log.Printf("Set() key:[%s] to %s \n", string(key), string(val))				//:todo is for testing only

	// time to live implementation:
	go func(){
		<-time.After(ttl)
		delete(c.data, string(key))		//:todo what happens if it is already deleted from ex a DELETE call?
	}()

	return nil	// :todo can this even fail? if not remove return

}


func (c *Cache) Get(key[]byte) ([]byte, error){
	c.lock.RLock()
	defer c.lock.RUnlock()
	log.Printf("Get() key:[%s] \n", string(key))								//:todo is for testing only
	val, ok := c.data[string(key)] 
	if ok{
		log.Printf("Get() key:[%s] result is %s \n", string(key), string(val))	//:todo is for testing only
		return val, nil
	}
	return nil, fmt.Errorf("key %s not found\n", string(key))
}	
	
func (c *Cache) Has(key[]byte) bool{
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Delete(key []byte) error{
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[string(key)]
	if ok{
		delete(c.data, string(key))
		return nil
	}
	return fmt.Errorf("Can't delete. key %s does not exist.\n", string(key))
}