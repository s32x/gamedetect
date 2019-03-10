package classifier /* import "s32x.com/tfclass/classifier" */

import (
	"errors"
	"sort"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

// Prediction is a struct containing a class label and the probability of the
// classified image being the cooresponding label
type Prediction struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

// ClassifyImage classifies a passed images bytes Buffer and returns the
// predictions as a slice of Predictions
func (c *Classifier) ClassifyImage(filename string, img []byte) ([]Prediction, error) {
	// Split out the filenames extension
	fn := strings.Split(filename, ".")
	if len(fn) < 2 {
		return nil, errors.New("Invalid filename passed")
	}
	ext := fn[len(fn)-1]

	// Create the scope and the input/output operations for normalization
	s := op.NewScope()
	in, out := c.NormalizeOutputs(s, ext)

	// Create the graph on which the scope operates on
	graph, err := s.Finalize()
	if err != nil {
		return nil, err
	}

	// Create a temporary session and defer close it
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer sess.Close()

	// Create a tensor from the stringified image buffer
	tensor, err := tf.NewTensor(string(img))
	if err != nil {
		return nil, err
	}

	// Normalize the temporary tensor and return the processed Classification
	result, err := sess.Run(map[tf.Output]*tf.Tensor{in: tensor},
		[]tf.Output{out}, nil)
	if err != nil {
		return nil, err
	}
	return c.ClassifyTensor(result[0])
}

// ClassifyTensor processes the passed tensor using the Classifiers config
// and returns the perdictions as a slice of Predictions
func (c *Classifier) ClassifyTensor(tensor *tf.Tensor) ([]Prediction, error) {
	// Create the input and output operation values for processing
	in := c.graph.Operation(c.config.InputLayer).Output(0)
	out := c.graph.Operation(c.config.OutputLayer).Output(0)

	// Process the passed tensor using the classifiers session/graph and return
	// a mapped slice of sorted Predictions
	result, err := c.session.Run(map[tf.Output]*tf.Tensor{in: tensor},
		[]tf.Output{out}, nil)
	if err != nil {
		return nil, err
	}
	values := result[0].Value().([][]float32)[0]

	// Bind all prediction values to their corresponding labels and return
	return mapPredictions(c.labels, values)[:c.config.NumPredictions], nil
}

// NormalizeOutputs produces an input and output operation output for
// normalizing an image with the passex extension
func (c *Classifier) NormalizeOutputs(s *op.Scope, ext string) (tf.Output, tf.Output) {
	input := op.Placeholder(s, tf.String)

	// Decode the image
	var decode tf.Output
	switch ext {
	case "png":
		decode = op.DecodePng(s, input, op.DecodePngChannels(3))
	case "gif":
		decode = op.DecodeGif(s, input)
	case "bmp":
		decode = op.DecodeBmp(s, input, op.DecodeBmpChannels(3))
	default:
		decode = op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))
	}

	// Div and Sub perform (value-Mean)/Scale for each pixel
	return input, op.Div(s,
		op.Sub(s,
			// Resize to the configs length/width with bilinear interpolation
			op.ResizeBilinear(s,
				// Create a batch containing a single image
				op.ExpandDims(s,
					// Use decoded pixel values
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{c.config.Height, c.config.Width})),
			op.Const(s.SubScope("mean"), c.config.Mean)),
		op.Const(s.SubScope("scale"), c.config.Scale))
}

// mapPredictions takes a slice of labels and prediction values and returns
// them in a populated and sorted slice of Predictions
func mapPredictions(labels []string, values []float32) []Prediction {
	var p []Prediction
	for i, v := range values {
		p = append(p, Prediction{Label: labels[i], Probability: v})
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Probability > p[j].Probability
	})
	return p
}
