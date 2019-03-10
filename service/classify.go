package service /* import "s32x.com/tfclass/service" */

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"s32x.com/tfclass/classifier"
)

// Results is a struct containing the results of a classified image
type Results struct {
	Filename    string                  `json:"filename"`
	Predictions []classifier.Prediction `json:"predictions"`
	SpeedMS     int64                   `json:"speed_ms"`
}

// Classify is an echo Handler that processes an image and returns its
// predicted classifications
func (s *Service) Classify(c echo.Context) error {
	// Read the image from the request
	fh, err := c.FormFile("image")
	if err != nil {
		return newISErr(c, err)
	}

	// Open the image and defer close it
	file, err := fh.Open()
	if err != nil {
		return newISErr(c, err)
	}
	defer file.Close()

	// Copy the bytes into a new bytes Buffer, classify the buffered bytes and
	// return the calculated predictions
	var buf bytes.Buffer
	io.Copy(&buf, file)
	start := nowMS()
	predictions, err := s.classifier.ClassifyImage(fh.Filename, buf.Bytes())
	if err != nil {
		return newISErr(c, err)
	}
	speedMS := nowMS() - start
	return c.JSON(http.StatusOK, &Results{
		Filename:    fh.Filename,
		Predictions: predictions,
		SpeedMS:     speedMS,
	})
}

// newISErr takes an error and encodes it in a map as a basic JSON response
func newISErr(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError,
		map[string]string{"error": err.Error()})
}

// nowMS returns the current time in milliseconds
func nowMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
