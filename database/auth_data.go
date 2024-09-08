package database

import (
	"fmt"

	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
	"github.com/oomph-ac/api/internal"
	"github.com/surrealdb/surrealdb.go"
	"golang.org/x/exp/maps"
)

const (
	statementObtainAuth = "SELECT * FROM oomphAuth WHERE key=$key;"
)

// ObtainAuth sends a request to obtain authentication data matching the given key.
func ObtainAuth(key string) (types.DBAuthData, *errors.APIError) {
	// Create a new goroutine where we query the database to check for the authentication data.
	res, err := RunJob(func() interface{} {
		keys := internal.InfoPool.Get().(map[string]any)
		defer internal.InfoPool.Put(keys)

		maps.Clear(keys)
		keys["key"] = key

		// Query the database, and if for some reason we are unable to, return an error.
		dbRes, err := DB.Query(statementObtainAuth, keys)
		if err != nil {
			return errors.New(
				errors.APIDatabaseFailed,
				"failed to query database for auth",
				err,
			)
		}

		// Parse the results found in the database query.
		var results []types.DBAuthData
		found, err := surrealdb.UnmarshalRaw(dbRes, &results)
		if err != nil {
			return errors.New(
				errors.APIDatabaseFailed,
				"cannot parse auth response from database",
				err,
			)
		}

		if !found {
			return errors.New(
				errors.APIUserFault,
				fmt.Sprintf("invalid authentication key (%s)", key),
				nil,
			)
		}

		return results[0]
	})

	if err != nil {
		return types.DBAuthData{}, err
	}
	return res.(types.DBAuthData), nil
}
