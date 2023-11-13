package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
)

type CSV struct {
	Path string `json:"path"`
}

func NewCSV(path string) *CSV {
	return &CSV{path}
}

func (c *CSV) Import() ([]*entity.Account, error) {
	if strings.TrimSpace(c.Path) == "" {
		return nil, errors.New("path must not be empty")
	}
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, errors.New("file not found: " + c.Path)
	}
	defer func() {
		file.Close()
	}()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var accounts []*entity.Account
	for _, r := range records[1:] {
		balance, err := strconv.ParseFloat(r[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to int: %v", err)
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
