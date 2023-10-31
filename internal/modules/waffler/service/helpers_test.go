//go:build unit
// +build unit

package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckOrdersParams(t *testing.T) {
	testData := []struct {
		err    error
		orders []string
	}{
		{err: nil, orders: []string{}},
		{err: nil, orders: []string{"test"}},
		{err: nil, orders: []string{"test1", "test2"}},
		{err: nil, orders: []string{"test1", "test2_desc"}},
		{err: nil, orders: []string{"test1", "test2", "test3"}},
		{err: nil, orders: []string{"test3_desc", "test2", "test1_desc"}},
		{err: badOrderParams, orders: []string{"test", "test_desc"}},
		{err: badOrderParams, orders: []string{"test_desc", "test"}},
		{err: badOrderParams, orders: []string{"test1", "test2", "test1_desc"}},
		{err: badOrderParams, orders: []string{"test1_desc", "test2", "test1"}},
		{err: badOrderParams, orders: []string{"test1", "test2", "test3_desc", "test1"}},
	}

	for count, tt := range testData {
		t.Run(fmt.Sprint("test ", count), func(t *testing.T) {
			err := checkOrderParams(tt.orders)
			assert.Equal(t, tt.err, err)
		})
	}
}
