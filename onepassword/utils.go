package onepassword

import (
	"strings"
)

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
