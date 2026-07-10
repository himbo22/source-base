package dto

// ZMember represents a member of a sorted set with its score.
type ZMember struct {
	Score  float64
	Member any
}
