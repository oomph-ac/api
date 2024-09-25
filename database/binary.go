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
	statementFindBinary = "SELECT data FROM oomph_bins WHERE (os=$os && arch=$arch && branch=$branch);"
)

// UpdateBinary updates the current Oomph binary data in the database.
func UpdateBinary(os, arch, branch, data string) *errors.APIError {
	return nil
}

// SearchForBinary searches for a binary in the database and returns the data if the
// specified binary type is found. An error is returned if the binary could not be found.
// We always want to send the latest binary in the database, don't cache results.
func SearchForBinary(os, arch, branch string) (string, *errors.APIError) {
	keys := internal.InfoPool.Get().(map[string]any)
	defer internal.InfoPool.Put(keys)

	maps.Clear(keys)
	keys["os"], keys["arch"], keys["branch"] = os, arch, branch

	dbRes, dbErr := DB.Query(statementFindBinary, keys)
	if dbErr != nil {
		return "", errors.New(
			errors.APIDatabaseFailed,
			"database query failed",
			dbErr,
		)
	}

	var results []types.DBProxyBinaryResponse
	if found, err := surrealdb.UnmarshalRaw(dbRes, &results); err != nil {
		return "", errors.New(
			errors.APIDatabaseFailed,
			"unable to parse response from database",
			err,
		)
	} else if !found {
		return "", errors.New(
			errors.APIUserFaultNeedsLog,
			fmt.Sprintf("could not find binary for binary %s_%s_%s", os, arch, branch),
			nil,
		)
	}

	return results[0].Data, nil
}
