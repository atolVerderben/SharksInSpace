package spaceshark

type Faction int

const (
	FactionPlayer Faction = iota
	FactionFriendly
	FactionNeutral
	FactionHostile
)
