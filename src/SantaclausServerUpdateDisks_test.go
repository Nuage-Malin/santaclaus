package main

// todo put this file in different directory

import (
	"testing"
)

/* AddFile */

func TestUpdateDisks(t *testing.T) {
	r := server.updateDiskBase(server.ctx)
	if r != nil {
		t.Error(r)
	}

}
