package schemas

// type ResultSchema struct {
// 	ScanStatus scanstatus.ScanStatus
// 	Success    bool
// 	Message    []string
// 	Data       []interface{}
// }

type InputItem struct {
	Target string
	Force  bool
}
