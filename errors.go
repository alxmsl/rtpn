package cpn

import "errors"

var ErrorEntityAlreadyExists = errors.New("entity already exists")
var ErrorEntityNotFound = errors.New("entity not found")
var ErrorNetIsActive = errors.New("net is active")
var ErrorNetIsInactive = errors.New("net is inactive")
var ErrorWrongTokenType = errors.New("wrong token type")
