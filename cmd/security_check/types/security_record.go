package types

type SecurityRecord struct {
	Fatal   bool
	Module  string
	Content string
	Guide   string
}

func NewSecurityRecord(module string, fatal bool, content string) SecurityRecord {
	return SecurityRecord{
		Fatal:   fatal,
		Module:  module,
		Content: content,
	}
}

func (r SecurityRecord) WithGuide(guide string) SecurityRecord {
	r.Guide = guide
	return r
}

func (r SecurityRecord) String() string {
	var result string
	if r.Fatal {
		result += "*FATAL*"
	} else {
		result += "WARNING"
	}
	result += " [" + r.Module + "] " + r.Content
	if r.Guide != "" {
		result += "\n> " + r.Guide
	}
	return result
}
