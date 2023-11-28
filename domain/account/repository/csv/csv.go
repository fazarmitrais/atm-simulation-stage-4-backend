package csv

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
)

type CSV struct {
}

func NewCSV() *CSV {
	return &CSV{}
}

func (c *CSV) Import(path string) ([]*entity.Account, *echo.HTTPError) {
	if strings.TrimSpace(path) == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "path must not be empty")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid file path: "+path)
	}
	defer func() {
		file.Close()
	}()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var accounts []*entity.Account
	for _, r := range records[1:] {
		balance, err := strconv.ParseFloat(r[2], 64)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to convert to int: %v", err))
		}
		accounts = append(accounts, &entity.Account{
			Name:          r[0],
			PIN:           r[1],
			Balance:       balance,
			AccountNumber: r[3],
		})
	}

	return accounts, nil
}
