package schemaversion

func V40(dash map[string]any) error {
	dash["schemaVersion"] = int(40)
	if _, ok := dash["refresh"].(string); !ok {
		dash["refresh"] = ""
	}
	return nil
}
