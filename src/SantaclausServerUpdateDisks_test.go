package main

// todo put this file in different directory

import (
	"context"
	"testing"
	"time"
)

/* AddFile */

// TODO do this test !
func TestUpdateDisks(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	r := TestServer.updateDiskBase(ctx)

	if r == nil { // updateDiskBase fails because Bugle is not conected
		t.Fatalf("updateDiskBase should not work if Bugle is not launched")
	}
}
