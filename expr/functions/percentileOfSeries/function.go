package percentileOfSeries

import (
	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
)

func init() {
	f := &percentileOfSeries{}
	functions := []string{"percentileOfSeries"}
	for _, function := range functions {
		metadata.RegisterFunction(function, f)
	}
}

type percentileOfSeries struct {
	interfaces.FunctionBase
}

// percentileOfSeries(seriesList, n, interpolate=False)
func (f *percentileOfSeries) Do(e parser.Expr, from, until int32, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	// TODO(dgryski): make sure the arrays are all the same 'size'
	args, err := helper.GetSeriesArg(e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	percent, err := e.GetFloatArg(1)
	if err != nil {
		return nil, err
	}

	interpolate, err := e.GetBoolNamedOrPosArgDefault("interpolate", 2, false)
	if err != nil {
		return nil, err
	}

	return helper.AggregateSeries(e, args, func(values []float64) float64 {
		return helper.Percentile(values, percent, interpolate)
	})
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *percentileOfSeries) Description() map[string]*types.FunctionDescription {
	return map[string]*types.FunctionDescription{
		"percentileOfSeries": {
			Description: "percentileOfSeries returns a single series which is composed of the n-percentile\nvalues taken across a wildcard series at each point. Unless `interpolate` is\nset to True, percentile values are actual values contained in one of the\nsupplied series.",
			Function:    "percentileOfSeries(seriesList, n, interpolate=False)",
			Group:       "Combine",
			Module:      "graphite.render.functions",
			Name:        "percentileOfSeries",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "n",
					Required: true,
					Type:     types.Integer,
				},
				{
					Default: "false",
					Name:    "interpolate",
					Type:    types.Boolean,
				},
			},
		},
	}
}