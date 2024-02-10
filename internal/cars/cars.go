package cars

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Car struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Discount int64  `json:"discount"`
}

type UpdateCar struct {
	Name     *string `json:"name,omitempty"`
	Price    *int64  `json:"price,omitempty"`
	Discount *int64  `json:"discount,omitempty"`
}

func GetCars(c echo.Context) error {
	carFile, err := os.Open("cars.json")

	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "interal server error",
			"error":   err.Error(),
		})
	}

	defer carFile.Close()

	var cars []Car
	json.NewDecoder(carFile).Decode(&cars)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(cars)
}

func GetCar(c echo.Context) error {
	carFile, err := os.Open("cars.json")

	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "interal server error",
			"error":   err.Error(),
		})
	}

	defer carFile.Close()

	var cars []Car
	json.NewDecoder(carFile).Decode(&cars)

	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	var car Car
	found := false

	for i := range cars {
		if cars[i].ID == idInt {
			car = cars[i]
			found = true
			break
		}
		continue
	}

	if !found {
		return c.JSON(404, map[string]interface{}{
			"message": "car " + id + " not found",
		})
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(&car)
}

func CreateCar(c echo.Context) error {
	bytes, err := os.ReadFile("cars.json")
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "interal server error",
			"error":   err.Error(),
		})
	}

	var cars []Car
	json.Unmarshal(bytes, &cars)

	var newCar Car
	if err := c.Bind(&newCar); err != nil {
		c.Logger().Error(err)
		return c.JSON(400, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	id := len(cars) + 1
	newCar.ID = id
	cars = append(cars, newCar)

	carsBytes, _ := json.Marshal(cars)
	err = os.WriteFile("cars.json", carsBytes, os.ModePerm)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "add car failed",
			"error":   "write file error due to " + err.Error(),
		})
	}

	return c.JSON(201, map[string]interface{}{
		"message": "new car added",
		"data":    newCar,
	})
}

func EditCar(c echo.Context) error {
	var updatePayload UpdateCar
	if err := c.Bind(&updatePayload); err != nil {
		c.Logger().Error(err)
		return c.JSON(400, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	var cars []Car
	carsBytes, err := os.ReadFile("cars.json")
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "interal server error",
			"error":   err.Error(),
		})
	}

	json.Unmarshal(carsBytes, &cars)
	var targetCar *Car
	found := false
	for i := range cars {
		if cars[i].ID == idInt {
			found = true
			targetCar = &cars[i]
			break
		}
		continue
	}
	if !found {
		return c.JSON(404, map[string]interface{}{
			"message": "car " + id + " does not exist",
		})
	}

	if updatePayload.Discount != nil {
		targetCar.Discount = *updatePayload.Discount
	}
	if updatePayload.Name != nil {
		targetCar.Name = *updatePayload.Name
	}
	if updatePayload.Price != nil {
		targetCar.Price = *updatePayload.Price
	}

	updatedCarsBytes, _ := json.Marshal(cars)
	err = os.WriteFile("cars.json", updatedCarsBytes, os.ModePerm)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "update car failed",
			"error":   "write file error due to " + err.Error(),
		})
	}

	return c.JSON(201, map[string]interface{}{
		"message": "car updated",
		"data":    *targetCar,
	})
}
