package goptions

// Help Defines the common help flag. It is handled separately as it will cause
// Parse() to return ErrHelpRequest.
type Help bool

// Verbs marks the point in the struct where the verbs start.
type Verbs interface{}

// A remainder catches all excessive arguments.
type Remainder []string
