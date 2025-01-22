package cpn

import "errors"

// ErrorEntityAlreadyExists defines error for a case when one entity conflicts with another.
var ErrorEntityAlreadyExists = errors.New("entity already exists")

// ErrorEntityNotFound defines error for a case when given entity has not been found.
var ErrorEntityNotFound = errors.New("entity not found")

// ErrorNetIsActive defines error for a case when Net is active unexpectedly.
var ErrorNetIsActive = errors.New("net is active")

// ErrorNetIsInactive defines error for a case when Net is inactive unexpectedly.
var ErrorNetIsInactive = errors.New("net is inactive")

// ErrorWrongTokenType defines error for a case when Token has a wrong type.
var ErrorWrongTokenType = errors.New("wrong token type")
