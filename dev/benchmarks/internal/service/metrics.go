package service

import "math"

// positions is a slice of 1 indexed ranks for each query, shows the position ofthe
// eval dataset issue appears in the search query results.
// returns 0 if the eval test set issue never appeared in results

// how high relevant doc ranks on average
func MRR(positions []int) float64 {
	if len(positions) == 0 {
		return 0
	}

	sum := 0.0
	for _, r := range positions {
		if r > 0 {
			sum += 1.0 / float64(r)
		}
	}

	return sum / float64(len(positions))
}

// did the relevant doc appear anywhere in top 10
func RecallAt10(positions []int) float64 {
	if len(positions) == 0 {
		return 0
	}

	hits := 0
	for _, r := range positions {
		if r > 0 && r <= 10 {
			hits++
		}
	}

	return float64(hits) / float64(len(positions))
}

// similar to MRR but less aggresive about penalizing
// rank 1 vs rank 2
func NDCGAt10(positions []int) float64 {
	if len(positions) == 0 {
		return 0
	}

	sum := 0.0
	for _, r := range positions {
		if r > 0 && r <= 10 {
			sum += 1.0 / math.Log2(float64(r)+1)
		}
	}

	return sum / float64(len(positions))
}
