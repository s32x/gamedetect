package service /* import "s32x.com/gamedetect/service" */

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"s32x.com/gamedetect/classifier"
)

// Results is a struct containing the results of a classified image
type Results struct {
	Filename    string                  `json:"filename,omitempty"`
	Predictions []classifier.Prediction `json:"predictions,omitempty"`
	SpeedMS     int64                   `json:"speed_ms,omitempty"`
}

// Classify is an echo Handler that processes an image and returns its
// predicted classifications
func (s *Service) Classify(c echo.Context) error {
	// Decode the image from the request
	filename, bytes, err := decodeFile(c, "image")
	if err != nil {
		return newISErr(c, err)
	}

	// Perform the classification and return the results
	start := nowMS()
	predictions, err := s.classifier.ClassifyImage(filename, bytes)
	if err != nil {
		return newISErr(c, err)
	}
	return c.JSON(http.StatusOK, &Results{
		Filename:    filename,
		Predictions: predictions,
		SpeedMS:     nowMS() - start,
	})
}

// decodeFile decodes a file from the passed echo Contexts form and returns
// both the files name and it's bytes
func decodeFile(c echo.Context, name string) (string, []byte, error) {
	// Read the image from the request
	fh, err := c.FormFile(name)
	if err != nil {
		return "", nil, err
	}

	// Open the image and defer close it
	file, err := fh.Open()
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	// Read the bytes into a bytes buffer and return
	var buf bytes.Buffer
	io.Copy(&buf, file)
	return fh.Filename, buf.Bytes(), nil
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
