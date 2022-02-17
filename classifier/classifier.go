// This package is an extremely naive implementation of a tensorflow image
// classification wrapper. It abstracts away a good amount of the boilerplate
// required to load and process images using the model/label outputs from
// tensorhubs retrain.py
// See: https://github.com/tensorflow/hub/blob/master/examples/image_retraining/retrain.py

package classifier

import (
	"bufio"
	"io/ioutil"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

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
	return &Classifier{config: config, graph: graph, session: session,
		labels: labels}, nil
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
