package models

type TempletData struct {
	StringMap    map[string]string
	IntMap       map[string]int
	FlotMap      map[string]float32
	Data         map[string]interface{}
	CSRFToken    string
	InfoEdit     string
	WarningEdit  string
	CriticalEdit string
}