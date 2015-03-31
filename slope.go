// Package slopeone provides an implementation of the Slope One algorithm
// used for collaborative filtering.
//
// The algorithm is introduced in:
//	Slope One Predictors for Online Rating-Based Collaborative Filtering (2005)
//
//	Daniel Lemire and Anna Maclachlan
//
// This paper is available here: http://lemire.me/fr/documents/publications/lemiremaclachlan_sdm05.pdf
//
// Slope One is an incredibly simple item-item collaborative filtering
// algorithm, which uses user-item ratings to provide a model to predict
// users' ratings for items they have yet to rate.
package slopeone

// This type is just for semantic intent.

// UserRatings is a set of item ratings belonging to a user.
type UserRatings map[int]float64

// S1 implements the Slope One algorithm.
type S1 struct {
	// d maintains a mapping between items and their rating differences
	// to other items. For examples, given item1 with a rating of 3.5
	// and item2 with a rating of 4.5, one could add the following to
	// the d:
	//	d["item1"]["item2"] = -1.0
	d map[int]map[int]float64

	// f maintains a mapping between items and the number of times
	// differenes in ratings have been calculated for other items.
	// For example, if the difference between item1 and item2 was
	// calculated, then the following would be added to f:
	//	f["item1"]["item2"]++
	f map[int]map[int]int
}

// NewS1 returns an *S1 ready for use.
func NewS1() *S1 {
	return &S1{
		d: make(map[int]map[int]float64),
		f: make(map[int]map[int]int),
	}
}

// AddRatings adds user ratings for sets of items to the S1.
// Ratings for added items will be taken into consideration in future
// predictions.
func (s1 *S1) AddRatings(users []UserRatings) {
	for _, user := range users {
		// For each item and rating generate the difference in rating
		// between this one and all other items.
		for i1, r1 := range user {
			if _, ok := s1.d[i1]; !ok {
				s1.d[i1] = make(map[int]float64)
				s1.f[i1] = make(map[int]int)
			}

			// Update the frequency of i1 vs i2 and the total rating
			// difference observed.
			for i2, r2 := range user {
				s1.f[i1][i2]++
				s1.d[i1][i2] += (r1 - r2)
			}
		}
	}

	// Normalise the difference in ratings for each pair of items, by
	// the number of times each item-pair have had their differences
	// calculated.
	for i1, diffs := range s1.d {
		for i2 := range diffs {
			diffs[i2] /= float64(s1.f[i1][i2])
		}
	}
}

// Predict returns predicted ratings for items the provided user has not
// yet rated, based on the rating they provide for items they have
// rated.
//
// Items the user has rated are not included in the returned
// UserPredictions.
func (s1 *S1) Predict(ur UserRatings) map[int]float64 {
	p, f := make(map[int]float64), make(map[int]int)
	var gf int
	// For each item-rating the user has rated we will compare it to
	// all global item-ratings, and update our prediction of unrated
	// items for the user.
	for i, r := range ur {
		for gi, gr := range s1.d {
			// If items have never been analysed or we will want to
			// remove them from the predicted set anyway, then move on.
			if gf = s1.f[gi][i]; gf == 0 || gi == i {
				continue
			}

			// Update our prediction of the unrated item's rating for
			// the user according to the global rating difference
			// between the user's rated item (i) and the other item
			// we're looking at (gi). This difference gives us a
			// direction to modify they user's providing rating for i
			// by, in order to predict their rating of gi.
			p[gi] += (float64(gf) * (gr[i] + r))
			f[gi] += gf
		}
	}

	// Normalise each predicted rating, and remove ones that were in the
	// set of provided ratings.
	for i := range p {
		p[i] /= float64(f[i])
		for j := range ur {
			if i == j {
				delete(p, j)
			}
		}
	}
	return p
}
