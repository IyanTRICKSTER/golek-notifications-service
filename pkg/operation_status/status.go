package status

type OperationStatus int

var Success OperationStatus = 1
var Failed OperationStatus = -1
var ErrorDuplicatedModel OperationStatus = -2
var SendMessageSuccess OperationStatus = 100
var SendMessageFailed OperationStatus = 101
