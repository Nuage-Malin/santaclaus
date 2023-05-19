package main

// todo put this file in different directory

import (
	context "context"
	"testing"
	"time"
)

/* AddFile */

// TODO do this test !
func TestUpdateDisks(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	r := server.updateDiskBase(ctx)

	if r != nil {
		t.Fatal(r)
	}

}
