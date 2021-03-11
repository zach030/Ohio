package znet

import (
	"ohio/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

func (c *ConnManager) Add(connection ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[connection.GetConnID()] = connection
	fmt.Println("ConnId = ", connection.GetConnID(), "connection add to ConnManager successfully, conn num = ", c.Len())
}

func (c *ConnManager) Remove(connection ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, connection.GetConnID())
	fmt.Println("ConnId = ", connection.GetConnID(), " delete successfully")
}

func (c *ConnManager) Get(id uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[id]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) Clear() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for id, connection := range c.connections {
		connection.Stop()
		delete(c.connections, id)
	}
	fmt.Println("Clear all connections successfully!")
}
