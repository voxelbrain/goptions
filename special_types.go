package goptions

// Help Defines the common help flag. It is handled separately as it will cause
// Parse() to return ErrHelpRequest.
type Help bool

// Verbs marks the point in the struct where the verbs start.
type Verbs interface{}

type Remainder []string
