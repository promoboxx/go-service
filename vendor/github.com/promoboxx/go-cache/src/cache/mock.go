package cache

type mock struct{}

// NewMock returns a Cache that doesn't actually cache anything
func NewMock() Cache {
	return &mock{}
}

func (m *mock) GetAndLoad(key string, value interface{}, loader func() (interface{}, error)) error {
	result, err := loader()
	if err != nil {
		return err
	}

	SetValue(value, result)
	return nil
}
