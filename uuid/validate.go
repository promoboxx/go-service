package uuid

import "regexp"

// IsValid returns true if the uuid string provided is in a valid format
// Note that this function was copied from https://stackoverflow.com/questions/25051675/how-to-validate-uuid-v4-in-go
func IsValid(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
