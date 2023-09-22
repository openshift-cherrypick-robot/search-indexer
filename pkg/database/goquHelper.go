// Copyright Contributors to the Open Cluster Management project

package database

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"k8s.io/klog/v2"
)

func useGoqu(query string, params []interface{}) (q string, p []interface{}, er error) {
	dialect := goqu.Dialect("postgres")
	resources := goqu.S("search").Table("resources")
	edges := goqu.S("search").Table("edges")

	switch query {
	case "SELECT uid, data FROM search.resources WHERE cluster=$1":
		q, p, er = dialect.From(resources).Prepared(true).
			Select("uid", "data").Where(goqu.C("cluster").Eq(params[0])).ToSQL()

	case "INSERT into search.resources values($1,$2,$3) ON CONFLICT (uid) DO NOTHING":
		q, p, er = dialect.From(resources).Prepared(true).
			Insert().Rows(goqu.Record{"uid": params[0], "cluster": params[1], "data": params[2]}).
			OnConflict(goqu.DoNothing()).ToSQL()

	case "UPDATE search.resources SET data=$2 WHERE uid=$1":
		q, p, er = dialect.From(resources).Prepared(true).
			Update().Set(goqu.Record{"data": params[1].(string)}).Where(goqu.C("uid").Eq(params[0])).ToSQL()

	case "DELETE from search.resources WHERE uid IN ($1)":
		q, p, er = dialect.From(resources).Prepared(true).
			Delete().Where(goqu.C("uid").In(params)).ToSQL()

	case "DELETE from search.edges WHERE sourceid IN ($1) OR destid IN ($1)":
		q, p, er = dialect.From(edges).Prepared(true).
			Delete().Where(goqu.Or(goqu.C("sourceid").In(params), goqu.C("destid").In(params))).ToSQL()

	// Queries for EDGES table.
	case "SELECT sourceid, edgetype, destid FROM search.edges WHERE edgetype!='interCluster' AND cluster=$1":
		q, p, er = dialect.From(edges).Prepared(true).
			Select("sourceid", "edgetype", "destid").
			Where(goqu.C("edgetype").Neq("interCluster"), goqu.C("cluster").Eq(params[0])).ToSQL()

	case "INSERT into search.edges values($1,$2,$3,$4,$5,$6) ON CONFLICT (sourceid, destid, edgetype) DO NOTHING":
		q, p, er = dialect.From(edges).Prepared(true).
			Insert().Cols("sourceid", "sourcekind", "destid", "destkind", "edgetype", "cluster").Vals(params).
			OnConflict(goqu.DoNothing()).ToSQL()

	case "DELETE from search.edges WHERE sourceid=$1 AND destid=$2 AND edgetype=$3":
		q, p, er = dialect.From(edges).Prepared(true).
			Delete().Where(goqu.C("sourceid").Eq(params[0]), goqu.C("destid").Eq(params[1]), goqu.C("edgetype").
			Eq(params[2])).ToSQL()

	default:
		er = fmt.Errorf("Unable to build goqu query for [%s]", query)
	}

	if er != nil {
		klog.Errorf("Error building goqu query. Error: %+v Query: %s", er, query)
	}

	return q, p, er
}
