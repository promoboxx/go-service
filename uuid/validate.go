package uuid

import googleuuid "github.com/google/uuid"

// IsValid returns true if the uuid string provided is in a valid format
// Note that this function was copied from https://stackoverflow.com/questions/25051675/how-to-validate-uuid-v4-in-go
func IsValid(uuid string) bool {
	var isValid bool
	_, err := googleuuid.Parse(uuid)
	if err == nil {
		isValid = true
	}
	return isValid
}
