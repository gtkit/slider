package slider

const (
	ErrSliderNotExists Error = "slider not exists"
	ErrSliderMismatch  Error = "slider mismatch"
	ErrSliderSave      Error = "slider save image error"
	ErrSliderVerify    Error = "slider verify error"
	ErrSliderReload    Error = "slider reload"
	ErrSliderCtxDone   Error = "slider context done"
	ErrSliderRefresh   Error = "slider refresh"
)

type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}
