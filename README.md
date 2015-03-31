# slopeone

[![GoDoc](https://godoc.org/github.com/e-dard/slopeone?status.svg)](https://godoc.org/github.com/e-dard/slopeone)

A Go implementation of the Slope One Collaborative Filtering algorithm, as introduced in:

> [Slope One Predictors for Online Rating-Based Collaborative Filtering](http://lemire.me/fr/documents/publications/lemiremaclachlan_sdm05.pdf)
>
> Daniel Lemire and Anna Maclachlan (2005)

## A What?

Collaborative filtering is an approach to implementing part or all of a recommendation system, wherein you recommend or suggest items to users by analysing other users' preferences to those items.

Typically, CF algorithms fall into two types:

#### User-based algorithms

These algorithms focus on identifying similarities between the user you want to predict item preferences for *u1*, and other uses in the system. Once you have identified other similar users you then use their preferences for items to predict *u1*'s preferences.

Slope One is not one of these algorithms

#### Item-based algorithms

Item-based algorithms don't consider the relationship between users at all, instead they use the relationship between *pairs of items*.
Slope One is one of these types of algorithms (maybe the simplest, actually).
The general idea is a sparse item-item matrix is generated, where each cell within the matrix describes a difference in preferences between those items.
These differences are calculated for each user in the dataset, but they're merged to give a global matrix describing how all users perceives differences between pairs of items.

To make a prediction for a new user, some existing preferences for that user need to be provided.
Using these preferences, predictions of the user's preferences of all other items in the system can be extracted.

## Usage

```
$ go get github.com/e-dard/slopeone
```

> Note: Currently it's not possible to add new user preferences to the
> model, or to update existing ones, though I might add this in the
> future. A basic use-case for the current model would be to
> recalculate all recommendations for all users, in batches, and then
> cache the results until the next time you run the model (with more
> data).

To incorporate the model, you can do something like this:

```go
package main

import (
	"fmt"

	"github.com/e-dard/slopeone"
)

func main() {
	// Each UserRatings map represents a single user's set of preferences
	// over some items. In this case we have three users' preferences.
	userRatings := []slopeone.UserRatings{
		{2005: 2.4, 5513: 1.3, 13035: 2.0},
		{5513: 4, 359602: 5, 13035: 1.5, 29074: 4},
		{29074: 4.3, 359602: 2.5, 2005: 5},
	}

	s1 := slopeone.NewS1()

	// Add user preferences to the model.
	s1.AddRatings(userRatings)

	// Predict will return predicted ratings for a user
	preds := s1.Predict(slopeone.UserRatings{2005: 2.0, 29074: 3.2})
	for item, rating := range preds {
		fmt.Printf("item: %d\trating: %0.1f\n", item, rating)
	}
}
```
results in:

```
item: 359602	rating: 1.7
item: 5513		rating: 2.1
item: 13035		rating: 1.2
```

