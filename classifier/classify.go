package classifier

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"mime/multipart"
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

// ClassifyImage takes an image Image, writes it to a bytes Buffer, performs a
// classification and returns a slice of predictions
func (c *Classifier) ClassifyImage(img image.Image) ([]Prediction, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return c.ClassifyBytes(buf.Bytes(), "png")
}

// ClassifyMultipart takes a multipart Fileheader, performs a classification
// and returns a slice of predictions
func (c *Classifier) ClassifyMultipart(fh *multipart.FileHeader) ([]Prediction, error) {
	// Open the FileHeader and defer close it
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the bytes into a bytes buffer and return
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return nil, err
	}

	// Split out the filenames extension and return a full classification on
	// the decoded bytes
	fn := strings.Split(fh.Filename, ".")
	if len(fn) < 2 {
		return nil, errors.New("Invalid filename passed")
	}
	return c.ClassifyBytes(buf.Bytes(), fn[len(fn)-1])
}

// ClassifyBytes classifies the passed images bytes and returns the predictions
// as a slice of Predictions
func (c *Classifier) ClassifyBytes(img []byte, ext string) ([]Prediction, error) {
	// Create the scope and the input/output operations for normalization
	s := op.NewScope()
	in, out := normalize(s, ext, c.config)

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
	return c.classifyTensor(result[0])
}

// classifyTensor processes the passed tensor using the Classifiers config
// and returns the perdictions as a slice of Predictions
func (c *Classifier) classifyTensor(tensor *tf.Tensor) ([]Prediction, error) {
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

// normalize produces an input and output operation output for normalizing an
// image with the passed extension and Config
func normalize(s *op.Scope, ext string, config Config) (tf.Output, tf.Output) {
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
	default: // Will default to Jpeg decode
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
				op.Const(s.SubScope("size"), []int32{config.Height, config.Width})),
			op.Const(s.SubScope("mean"), config.Mean)),
		op.Const(s.SubScope("scale"), config.Scale))
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
