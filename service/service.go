package service /* import "s32x.com/tfclass/service" */

import (
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"s32x.com/tfclass/classifier"
)

// Service is a struct that contains everything needed to perform image
// predictions
type Service struct {
	mu          sync.Mutex
	classifier  *classifier.Classifier
	testResults []Result
}

// NewService creates a new Service reference using the given service params
func NewService(graphPath, labelsPath string) (*Service, error) {
	// Create the game classifier using it's default config
	c, err := classifier.NewClassifier(graphPath, labelsPath)
	if err != nil {
		return nil, err
	}
	return &Service{classifier: c}, nil
}

// Close closes the Service by closing all it's closers ;)
func (s *Service) Close() error { return s.classifier.Close() }

// Start begins serving the generated Service on the passed port
func (s *Service) Start(port, testDir string) {
	// Process the testdata
	go s.TestData(testDir)

	// Create a new echo Echo and bind all middleware
	e := echo.New()
	e.HideBanner = true
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("service/templates/*.html")),
	}

	// Bind all middleware
	e.Pre(middleware.RemoveTrailingSlashWithConfig(
		middleware.TrailingSlashConfig{
			RedirectCode: http.StatusPermanentRedirect,
		}))
	e.Pre(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// Create the static file endpoints
	e.Static("*", "./service/static")

	// Bind all API endpoints
	e.POST("/", s.Classify)
	e.GET("/", s.Index)

	// Listen and Serve
	log.Printf("Starting service on port %v\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
