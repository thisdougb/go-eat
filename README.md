# go-eat

Welcome to Go Eat, my fictional restaurant.
Take a seat and we'll get started.

This is the journey from small app to scaled app.
A template of ideas that come from me doing this commercially.

I'm doing this because it's an easy way to prototype things without breaking my main app.
I can try out ideas, get feedback and discussion.
Hopefully some others will find it useful too.

## Operations

The restaurant is highly efficient.
Customers that love automation, love Go Eat.

We serve a set menu each day, no deviations.
This means our staff can focus on their job, not boring and repetitive admin.

### Booking

Input:

We take customer booking via a third-party app which calls our webhook endpoint.

Output:

After a successful reservation the kitchen is notified.

### Kitchen

Input:

Customer reservations schedule orders for the kitchen.
The kitchen stock system is automated, so outbound requests to suppliers may occur.

Output: 

When plates are ready a notification is sent to the service staff.

## The House Style

I love Go as a language.
Enough (type) safety without being pedantic.
And always falling on the side of simplicity is best.
Every language has pros and cons, so no more needs to be said.
This is all made with Go.

While I try to follow Go's idioms, sometimes I don't.
Overwhelmingly this is to do with pragmatism and that I'm not building simple packages.
So I step outside of the idiomatic Go layout in the interests of simplicity.

Every decision is made with an eye on future scalability.
I'm not building for >1,000,000 customers on day one.
But I am avoiding code that makes scaling complicated in the future.

## Monolith

The crucial point in scaling a business/app is when you hire enough developers to be able to specialise effort.
One, two, or five developers working on a shared code-base is usually fine, and easy to coordinate.
More than that and the code-base is too large to avoid each developer specialising.
At this point break your code into distinct repositories is the obvious choice.
And this is the point you have to completely refactor, or have an easy ride to a service architecture.

Go's package system is very simple, and really suits separated concerns in distinct repo's.
However, the practicality of continuously updating local copies becomes a real drag.
Deploying and versioning multiple services adds to the friction.

A monolith is simple.
Particularly with Go packages, a monolith layout can support a service architecture.
It is also very easy to deploy (and revert) a single process.

We will build and layout the monolith in a way that makes it easy to extract components.
If our kitchen component expands to make take-aways, we can refactor just that to an external repo and developers.
But it makes no sense to slow development today by externalising it.

A monolith because:

- this is intended to grow into a larger app 
- development begins with 1-4 developers who talk to each other
- low tolerance for admin overhead


### Package Layout

The initial package layout reflects the different domains/concerns within our app and business.
These could be apps or services in their own right.

In fact, when the business scales we will spin out packages as standalone services.
We may hire a dev team to work exclusively on the menu system.

This layout will make that transition easier, if it ever happens.
And in the mean time it keeps functionality separated in a Go-like way.

```
./goeat/booking
./goeat/menu
./goeat/staff
./goeat/service
./goeat/kitchen
./goeat/README.md
./goeat/main.go
./goeat/go.mod
```

### Domain Boundaries

One of the nasties that hits companies when the try and scale is _monolith mess_.
Code execution paths that are so intertwined that the pragmatic path to scale is rewriting the app into services.
That's a costly way of scaling.

I try to avoid creating this mess by using Go's `internal` package qualifier.
For example, when developing features in `kitchen` I want a definitive API into `staff`.

Go's package export method (capitalised type names) is too lax for my case.
I want to protect package exports inter-service, but use standard Go package exports intra-service.

As a short example, I add some files:

```
./goeat/kitchen/
./goeat/kitchen/internal/rota.go

./goeat/staff/
./goeat/staff/publicapi.go
./ // ./goeat/staff/publicapi.go
package staff

import "goeat/staff/internal/calendar"

var (
        GetKitchenRota = func() []string {
                return calendar.GetRota("kitchen")
        }
)
```

And calling this from the `kitchen` service looks like this:

``
package rota

import (
	"goeat/staff"
)

func fetchRota() []string {

	rota := staff.GetKitchenRota()

	// do something with the rota

	return rota
}
```

### Testing

Implementing the public 'api' methods as vars enables me to isolate testing between services.

It is trivial to mock the responses from the staff service, independently of that package.
I haven't found a simpler way of doing this, every other way seems to involve code gymnastics with interfaces or channels.
This particular use-case is emulating services within a monolith.

A simple test example in the kitchen service:

``` 
//go:build test

package rota

import (
	"goeat/staff"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchRota(t *testing.T) {

	mockRota := []string{"A", "B"}
	mockStaffGetKitchenRota := func() []string {
		return mockRota
	}

	// override the 'api' response with our mock function
	staff.GetKitchenRota = mockStaffGetKitchenRota

	assert.Equal(t, fetchRota(), mockRota, "expect A, B")
}
```

### Summary

My choice to develop an app as a monolith depends on a bunch of particular use-case factors.

I think Go's inherent simplicity around packages make it easy to create domain boundaries. 
This makes the whole code-base simpler to think about.

Simplicity is key, because with few (or solo) developers complexity kills motivation.
Speed of development is often less about CPU cycles, and more about efficient process.

A single deployable binary is easier to manage and test, at small scale.
Although Go suits a discrete service-architecture, the overhead (repo's, ci, infra) is too much right now.

