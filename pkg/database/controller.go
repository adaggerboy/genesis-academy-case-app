package database

import (
	"fmt"

	"github.com/adaggerboy/genesis-academy-case-app/models/config"
)

type DatabaseFabricFunc[R any] func(endpoint config.DatabaseEndpointConfig) (db R, err error)

type DatabaseFabric[R any] struct {
	readDBFabrics map[string]DatabaseFabricFunc[R]
}

func NewDatabaseFabric[R any]() *DatabaseFabric[R] {
	return &DatabaseFabric[R]{
		readDBFabrics: map[string]DatabaseFabricFunc[R]{},
	}
}

func (c *DatabaseFabric[R]) RegisterDatabaseFabric(key string, fun DatabaseFabricFunc[R]) {
	c.readDBFabrics[key] = fun
}

func (c *DatabaseFabric[R]) NewDatabaseController(conf config.DatabaseEndpointConfig) (controller R, err error) {

	err = nil

	readDBFabric, ok := c.readDBFabrics[conf.Driver]
	if !ok {
		err = fmt.Errorf("driver not found: %s", conf.Driver)
		return
	}
	controller, err = readDBFabric(conf)
	if err != nil {
		err = fmt.Errorf("read db controller init: %s", err)
		return
	}
	return
}
