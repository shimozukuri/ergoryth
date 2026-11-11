package lead

type Status string

const (
	StatusNew Status = "new"
)

func (s Status) IsValid() bool {
	return s == StatusNew
}
