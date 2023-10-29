package weather_api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/siri-aws-web-app/verdandi-weather-service/internal/database"
)

func StartWeatherApi() {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "Fiber",
		AppName:       "verdandi-weather-service",
	})

	app.Use(cors.New())

	app.Get("/current-weather-data", func(c *fiber.Ctx) error {
		cities, err := GetCitiesList(c.Query("cities"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		data, err := GetCurrentWeatherDataFromDb(cities)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get current weather data")
		}

		return c.SendString(string(data))
	})

	app.Get("/forecast-data", func(c *fiber.Ctx) error {
		cities, err := GetCitiesList(c.Query("cities"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		data, err := ForecastData(cities)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get forecast data")
		}

		return c.SendString(string(data))
	})

	app.Listen(":3000")
}

func GetCurrentWeatherDataFromDb(citiesList []string) ([]byte, error) {
	return database.GetCurrentWeatherDataFromDb(citiesList)
}

func ForecastData(citiesList []string) ([]byte, error) {
	return database.GetForecastDataFromDb(citiesList)
}

func GetCitiesList(cities string) ([]string, error) {
	if cities == "" {
		return nil, errors.New("cities parameter is required")
	}
	citiesList := strings.Split(cities, ",")
	return citiesList, nil
}
