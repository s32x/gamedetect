package service

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/s32x/gamedetect/classifier"
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
	// Read the FileHeader from the request
	fh, err := c.FormFile("image")
	if err != nil {
		return newISErr(c, err)
	}

	// Perform the classification and return the results
	start := nowMS()
	predictions, err := s.classifier.ClassifyMultipart(fh)
	if err != nil {
		return newISErr(c, err)
	}
	return c.JSON(http.StatusOK, &Results{
		Filename:    fh.Filename,
		Predictions: predictions,
		SpeedMS:     nowMS() - start,
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
