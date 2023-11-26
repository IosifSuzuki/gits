package dto

type Role string

const (
	Unknown Role = "unknown"
	Admin   Role = "admin"
	Viewer  Role = "viewer"
)

var allowedValues = []string{
	string(Admin),
	string(Viewer),
}

func ParseString(role string) Role {
	for _, value := range allowedValues {
		if role == value {
			return Role(role)
		}
	}
	return Unknown
}
