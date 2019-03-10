package service

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"s32x.com/tfclass/classifier"
)

// Result is a struct containing the results of a classified image
type Result struct {
	Filename    string                  `json:"filename"`
	Predictions []classifier.Prediction `json:"predictions"`
	SpeedMS     int64                   `json:"speed_ms"`
}

// Index is an echo Handler that renders the index template with the given data
func (s *Service) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"tests": s.testResults,
	})
}

// TestData pre-processes all the test data in the passed testDir, storing the
// results in the testResults map for observing
func (s *Service) TestData(testDir string) error {
	return filepath.Walk(testDir, func(path string, fi os.FileInfo, err error) error {
		// If there's an error reading the file or it's a a directory, return
		if err != nil || fi.IsDir() {
			return nil
		}
		filename := fi.Name()

		// Read in the files bytes
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Classify the image and calculate the speed of the classification in
		// milliseconds
		start := nowMS()
		predictions, err := s.classifier.ClassifyImage(filename, bytes)
		if err != nil {
			return err
		}
		speedMS := nowMS() - start

		// Store the test result in the testResults
		s.mu.Lock()
		s.testResults = append(s.testResults,
			Result{
				Filename:    filename,
				Predictions: predictions,
				SpeedMS:     speedMS,
			})
		s.mu.Unlock()
		return nil
	})
}
