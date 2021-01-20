package json

type PatchOp string

var (
	Add     PatchOp = "add"
	Replace PatchOp = "replace"
	Remove  PatchOp = "remove"

	// Request  JsonPatchType = "request"
	// Response JsonPatchType = "response"
)

type Patch struct {
	// Type  JsonPatchType `json:"-"`
	Op    PatchOp     `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
