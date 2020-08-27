package profile

func (a *AsIs) UnmarshalJSON(b []byte) error {
	*a = AsIs(string(b))
	return nil
}
