//go:build wireinject

package injector

import (
	"reflect"

	"github.com/google/wire"
	"github.com/umardev500/common/model"
	"github.com/umardev500/gochat/internal/routes"
)

type AllRoutes struct {
	Chat *routes.ChatRouteImpl `wire:"chatRoute"`
}

// 🔹 Provide Routes (Combines Chat Route)
func ProvideRoutes(allRoutes AllRoutes) []model.Route {
	var routes []model.Route

	v := reflect.ValueOf(allRoutes)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i).Interface()

		// Ensure the field is of type model.Route before adding
		if route, ok := fieldValue.(model.Route); ok {
			routes = append(routes, route)
		}
	}

	return routes
}

var RoutesSet = wire.NewSet(
	ChatRouteSet,
	wire.Struct(new(AllRoutes), "*"),
	ProvideRoutes,
)

// 🔹 Initialize Routes
func InitializeRoutes() []model.Route {
	wire.Build(
		RoutesSet,
	)

	return nil
}
