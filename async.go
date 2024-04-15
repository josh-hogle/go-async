package async

import "context"

// Future interface has the method signature for await.
type Future interface {
	Await() any
}

// future implements the Future interface.
type future struct {
	await func(ctx context.Context) any
}

// Await is used to wait for an async function to complete.
func (f future) Await() any {
	return f.await(context.Background())
}

// Exec executes the async function.
func Exec(f func() any) Future {
	var result any
	c := make(chan struct{})
	go func() {
		defer close(c)
		result = f()
	}()
	return future{
		await: func(ctx context.Context) any {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return result
			}
		},
	}
}
