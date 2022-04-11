// Copyright Contributors to the Open Cluster Management project
package clustersync

import (
	"context"
	"testing"
	"time"

	fakeClient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func Test_runLeaderElectionAndCancel(t *testing.T) {
	supressConsoleOutput()

	// Mock Kube client
	mockClient := fakeClient.NewSimpleClientset()
	lock := resourcelock.LeaseLock{
		Client: mockClient.CoordinationV1(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Mock the leader function.
	var started, cancelled bool
	mockStartLeaderFn := func(c context.Context) {
		started = true
		<-c.Done()
		cancelled = true
	}

	// Execute the function.
	go runLeaderElection(ctx, &lock, mockStartLeaderFn)

	// Wait for function to start and get cancelled.
	time.Sleep(1 * time.Millisecond)
	cancel()
	time.Sleep(1 * time.Millisecond)

	// Validate that the leader function was started and cancelled as expected.
	if !started {
		t.Error("Expected leader function to be executed.")
	}
	if !cancelled {
		t.Error("Expected leader process to be cancelled.")
	}
}
