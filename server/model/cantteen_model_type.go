package model

type ModelType int

// ModelType is used to differentiate between the type of the model for different queries to reduce duplicated code.
const (
	DISH      ModelType = 1
	CAFETERIA ModelType = 2
	NAME      ModelType = 3
)
