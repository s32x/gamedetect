// This package is an extremely naive implementation of a tensorflow image
// classification wrapper. It abstracts away a good amount of the boilerplate
// required to load and process images using the model/label output from
// tensorhubs retrain.py
// See: https://github.com/tensorflow/hub/blob/master/examples/image_retraining/retrain.py

package classifier /* import "s32x.com/gamedetect/classifier" */

import (
	"bufio"
	"io/ioutil"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

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

// Classifier is a struct used for classifying images
type Classifier struct {
	config  Config
	graph   *tf.Graph
	session *tf.Session
	labels  []string
}

// NewClassifier creates a new Classifier using the default config
func NewClassifier(graphPath, labelPath string) (*Classifier, error) {
	return NewClassifierWithConfig(graphPath, labelPath, DefaultConfig)
}

// NewClassifierWithConfig creates a new image Classifier for processing image
// predictions
func NewClassifierWithConfig(graphPath, labelPath string,
	config Config) (*Classifier, error) {
	// Read the passed inception model file
	bytes, err := ioutil.ReadFile(graphPath)
	if err != nil {
		return nil, err
	}

	// Populate a new graph using the read model
	graph := tf.NewGraph()
	if err := graph.Import(bytes, ""); err != nil {
		return nil, err
	}

	// Create a new execution session using the graph
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	// Read all labels in the passed labelPath
	labels, err := readLabels(labelPath)
	if err != nil {
		return nil, err
	}

	// Return a fully populated Classifier
	return &Classifier{config, graph, session, labels}, nil
}

// Close closes the Classifier by closing all it's closers ;)
func (c *Classifier) Close() error { return c.session.Close() }

// readLabels takes a path to a labels file, parses out, and returns all labels
// as a slice of strings
func readLabels(labelsPath string) ([]string, error) {
	// Open the passed labels file and defer close it
	f, err := os.Open(labelsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Scan all lines and populate a slice of labels
	var labels []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		labels = append(labels, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return labels, nil
}
