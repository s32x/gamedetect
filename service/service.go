package service /* import "s32x.com/gamedetect/service" */

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"s32x.com/gamedetect/classifier"
)

// Service is a struct that contains everything needed to perform image
// predictions
type Service struct{ classifier *classifier.Classifier }

// NewService creates a new Service reference using the given service params
func NewService(graphPath, labelsPath string) (*Service, error) {
	// Create the game classifier using it's default config
	c, err := classifier.NewClassifier(graphPath, labelsPath)
	if err != nil {
		return nil, err
	}
	return &Service{c}, nil
}

// Close closes the Service by closing all it's closers ;)
func (s *Service) Close() error { return s.classifier.Close() }

// Start begins serving the generated Service on the passed port
func (s *Service) Start(env, domain, demo, port string) {
	// Create a new echo Echo and bind all middleware
	e := echo.New()
	e.HideBanner = true

	// Configure SSL, WWW, and Host based redirects if being hosted in a
	// production environment
	if strings.Contains(strings.ToLower(env), "prod") {
		e.Pre(middleware.HTTPSNonWWWRedirect())
		e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if c.Request().Host == domain {
					return next(c)
				}
				return c.Redirect(http.StatusPermanentRedirect,
					c.Scheme()+"://"+domain)
			}
		})
		e.Pre(middleware.CORS())
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

	// If hosting as a demonstration, configure a renderer, process all demo
	// test data and serve all testdata and static assets on the index
	if strings.Contains(strings.ToLower(demo), "true") {
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
	}

	// Bind all API endpoints
	e.POST("/", s.Classify)

	// Listen and Serve
	log.Printf("Starting service on port %v\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
