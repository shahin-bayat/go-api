package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("John", "Doe", "password")
	if err != nil {
		assert.Nil(t, err)
	}

	fmt.Printf("Account: %+v\n", acc)

}
