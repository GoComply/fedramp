package validation_root

func (bm *BackMatter) GetResourceByUuid(uuid string) *Resource {
	for _, res := range bm.Resources {
		if res.Uuid == uuid {
			return &res
		}
	}
	return nil
}
