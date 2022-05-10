package bucket

type Bucket struct {
	label string
}

func New(label string) Bucket {
	return Bucket{label}
}

func (b Bucket) Label() string {
	return b.label
}

func (b Bucket) Weight() uint32 {
	// all buckets have even weight
	return 1
}
