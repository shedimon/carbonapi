package nonNegativeDerivative

import (
	"errors"
	"fmt"
	"github.com/bookingcom/carbonapi/expr/helper"
	"github.com/bookingcom/carbonapi/expr/interfaces"
	"github.com/bookingcom/carbonapi/expr/types"
	"github.com/bookingcom/carbonapi/pkg/parser"
	"math"
)

type nonNegativeDerivative struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &nonNegativeDerivative{}
	functions := []string{"nonNegativeDerivative"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

func (f *nonNegativeDerivative) Do(e parser.Expr, from, until int32, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	args, err := helper.GetSeriesArg(e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	maxValue, err := e.GetFloatNamedOrPosArgDefault("maxValue", 1, math.NaN())
	if err != nil {
		return nil, err
	}
	minValue, err := e.GetFloatNamedOrPosArgDefault("minValue", 2, math.NaN())
	if err != nil {
		return nil, err
	}
	hasMax := !math.IsNaN(maxValue)
	hasMin := !math.IsNaN(minValue)

	if hasMax && hasMin && maxValue <= minValue {
		return nil, errors.New("minValue must be lower than maxValue")
	}
	if hasMax && !hasMin {
		minValue = 0
	}

	argMask := 0
	if _, ok := e.NamedArgs()["maxValue"]; ok || len(e.Args()) > 1 {
		argMask |= 1
	}
	if _, ok := e.NamedArgs()["minValue"]; ok || len(e.Args()) > 2 {
		argMask |= 2
	}

	var result []*types.MetricData
	for _, a := range args {
		var name string
		switch argMask {
		case 3:
			name = fmt.Sprintf("nonNegativeDerivative(%s,%g,%g)", a.Name, maxValue, minValue)
		case 2:
			name = fmt.Sprintf("nonNegativeDerivative(%s,minValue=%g)", a.Name, minValue)
		case 1:
			name = fmt.Sprintf("nonNegativeDerivative(%s,%g)", a.Name, maxValue)
		case 0:
			name = fmt.Sprintf("nonNegativeDerivative(%s)", a.Name)
		}

		r := *a
		r.Name = name
		r.Values = make([]float64, len(a.Values))
		r.IsAbsent = make([]bool, len(a.Values))

		prev := a.Values[0]
		for i, v := range a.Values {
			if i == 0 || a.IsAbsent[i] || a.IsAbsent[i-1] {
				r.IsAbsent[i] = true
				prev = v
				continue
			}
			diff := v - prev
			if diff >= 0 {
				r.Values[i] = diff
			} else if hasMax && maxValue >= v {
				r.Values[i] = ((maxValue - prev) + (v - minValue) + 1)
			} else if hasMin && minValue <= v {
				r.Values[i] = (v - minValue)
			} else {
				r.Values[i] = 0
				r.IsAbsent[i] = true
			}
			prev = v
		}
		result = append(result, &r)
	}
	return result, nil
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *nonNegativeDerivative) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"nonNegativeDerivative": {
			Description: "Same as the derivative function above, but ignores datapoints that trend\ndown.  Useful for counters that increase for a long time, then wrap or\nreset. (Such as if a network interface is destroyed and recreated by unloading\nand re-loading a kernel module, common with USB / WiFi cards.\n\nExample:\n\n.. code-block:: none\n\n  &target=nonNegativederivative(company.server.application01.ifconfig.TXPackets)",
			Function:    "nonNegativeDerivative(seriesList, maxValue=None)",
			Group:       "Transform",
			Module:      "graphite.render.functions",
			Name:        "nonNegativeDerivative",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name: "maxValue",
					Type: types.Float,
				},
				{
					Name: "minValue",
					Type: types.Float,
				},
			},
		},
	}
}
