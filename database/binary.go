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
	statementFindBinary = "SELECT data FROM oomph_bins WHERE (os=$os && arch=$arch);"
)

// SearchForBinary searches for a binary in the database and
func SearchForBinary(os, arch string) (string, *errors.APIError) {
	res, err := RunJob(func() interface{} {
		keys := internal.InfoPool.Get().(map[string]any)
		defer internal.InfoPool.Put(keys)

		maps.Clear(keys)
		keys["os"], keys["arch"] = os, arch

		dbRes, dbErr := DB.Query(statementFindBinary, keys)
		if dbErr != nil {
			return errors.New(
				errors.APIDatabaseFailed,
				"database query failed",
				dbErr,
			)
		}

		var results []types.DBProxyBinaryResponse
		if found, err := surrealdb.UnmarshalRaw(dbRes, &results); err != nil {
			return errors.New(
				errors.APIDatabaseFailed,
				"unable to parse response from database",
				err,
			)
		} else if !found {
			return errors.New(
				errors.APIUserFaultNeedsLog,
				fmt.Sprintf("could not find binary for binary %s_%s", os, arch),
				nil,
			)
		}

		return results[0].Data
	})

	if err != nil {
		return "", err
	}
	return res.(string), nil
}
