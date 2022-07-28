package models

type Right uint32

const (
	Public Right = iota
	HasAccess
	GroupMember
	GroupACLKeeper
	GroupGateKeeper
	GroupOwner
	SBOwner
	Private
)
