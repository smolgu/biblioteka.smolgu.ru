package users

type Status int64

const (
	Banned Status = 1 << iota
	EmailConfirmed
	PersonConfirmed
	Admin

	StatusHolder5
	StatusHolder6
	StatusHolder7
	StatusHolder8
	StatusHolder9
	StatusHolder10
)

func (s Status) Add(toAdd Status) Status {
	return s | toAdd
}

func (s Status) Remove(toRemove Status) Status {
	if s&toRemove == 0 {
		return s
	}
	return s ^ toRemove
}

func (u User) IsAdmin() bool {
	return u.Status&Admin != 0
}
