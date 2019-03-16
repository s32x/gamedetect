package service

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo"
	"s32x.com/gamedetect/classifier"
)

// Template contains pre-processed templates for rendering
type Template struct{ ts *template.Template }

// Render satisfies the Renderer interface in order to populate a template page
func (t *Template) Render(w io.Writer, name string, data interface{},
	c echo.Context) error {
	return t.ts.ExecuteTemplate(w, name+".html", data)
}

// TestResults contains all test results and summary data
type TestResults struct {
	Correct  int
	Accuracy float64
	Results  []TestResult
}

// TestResult is a struct containing the results of a tested classified image
type TestResult struct {
	Filename    string                  `json:"filename,omitempty"`
	Expected    string                  `json:"expected,omitempty"`
	Correct     bool                    `json:"correct,omitempty"`
	Predictions []classifier.Prediction `json:"predictions,omitempty"`
	SpeedMS     int64                   `json:"speed_ms,omitempty"`
}

// Index is an echo Handler that renders the index template with the given data
func Index(tr *TestResults) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"tests": tr,
		})
	}
}

// ProcessTestData pre-processes all the test data in the passed testDir, returning
// the results for serving on the demos index
func ProcessTestData(classifier *classifier.Classifier, testDir string) (*TestResults, error) {
	tr := &TestResults{}
	if err := filepath.Walk(testDir, func(path string, fi os.FileInfo, err error) error {
		// If there's an error reading the file or it's a directory, return
		filename := fi.Name()
		if err != nil || fi.IsDir() || filename == ".DS_Store" {
			return nil
		}

		log.Println("Performing test on file:", path)

		// Get the expected value for the test data
		expected := strings.Split(strings.Replace(path, testDir, "", -1),
			string(os.PathSeparator))[1]

		// Read the files bytes
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Classify the image and calculate the speed of the classification in
		// milliseconds
		start := nowMS()
		predictions, err := classifier.ClassifyImage(filename, bytes)
		if err != nil {
			return err
		}
		speedMS := nowMS() - start

		// Store the test result in the slice of TestResults
		correct := expected == predictions[0].Label
		tr.Results = append(tr.Results, TestResult{
			Filename:    filename,
			Expected:    expected,
			Correct:     correct,
			Predictions: predictions,
			SpeedMS:     speedMS,
		})
		if correct {
			tr.Correct++
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Calculate and set overall accuracy
	var accuracy float64
	for _, res := range tr.Results {
		for _, p := range res.Predictions {
			if p.Label == res.Expected {
				accuracy = accuracy + float64(p.Probability)
			}
		}
	}
	tr.Accuracy = accuracy / float64(len(tr.Results))
	return tr, nil
}
