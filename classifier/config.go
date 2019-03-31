package classifier /* import "s32x.com/gamedetect/classifier" */

// Config is a struct used for configuring the classifier
type Config struct {
	Height, Width           int32
	Mean, Scale             float32
	InputLayer, OutputLayer string
	NumPredictions          int
}

// DefaultConfig is used in the case where a config is not defined by the user
var DefaultConfig = Config{
	Height:         299,
	Width:          299,
	Mean:           0,
	Scale:          255,
	InputLayer:     "Placeholder",
	OutputLayer:    "final_result",
	NumPredictions: 5,
}
