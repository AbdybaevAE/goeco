package codes

// Typed operation code(includes typed error that could appear)
type Code string

// Some known operation codes
const (
	Ok                 Code = "Ok"
	BadArguments       Code = "BadArguments"
	ServiceInternal    Code = "ServiceInternal"
	ServiceUnavailable Code = "ServiceUnavailable"
)
