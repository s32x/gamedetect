package service

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/s32x/gamedetect/classifier"
)

// Service is a struct that contains everything needed to perform image
// predictions
type Service struct {
	classifier *classifier.Classifier
	domain     string
	demo       string
}

// New creates a new Service reference using the given service params
func New(graphPath, labelsPath, domain, demo string) (*Service, error) {
	// Create the game classifier using it's default config
	c, err := classifier.NewClassifier(graphPath, labelsPath)
	if err != nil {
		return nil, err
	}
	return &Service{
		classifier: c,
		domain:     domain,
		demo:       demo,
	}, nil
}

// Close closes the Service by closing all it's closers ;)
func (s *Service) Close() error { return s.classifier.Close() }

// Echo creates and returns an echo router that serves the gamedetect api
func (s *Service) Echo() *echo.Echo {
	// Create a new echo Echo and bind all middleware
	e := echo.New()
	e.HideBanner = true

	// Bind middleware
	e.Pre(middleware.RemoveTrailingSlashWithConfig(
		middleware.TrailingSlashConfig{
			RedirectCode: http.StatusPermanentRedirect,
		}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	// If hosting as a demonstration, configure a renderer, process all demo
	// test data and serve all testdata and static assets on the index
	if strings.Contains(strings.ToLower(s.demo), "true") {
		e.Renderer = &Template{
			ts: template.Must(template.ParseGlob("service/templates/*.html")),
		}

		// Process the test data for serving in the index
		testData, err := ProcessTestData(s.classifier, "service/static/test")
		if err != nil {
			log.Fatal(err)
		}

		// Create the static file endpoints
		e.Static("*", "service/static")
		e.GET("/", Index(testData))
		e.POST("/demo", Demo(testData, s.classifier))
	}

	// Bind all API endpoints
	e.POST("/", s.Classify)
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	return e
}
