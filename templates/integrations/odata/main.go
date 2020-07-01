package odata

import "[[.Config.LocalProjectPath]]/types"

var (
	odataConfig types.OdataConfig
)

func SetOdataConfig(config types.OdataConfig) {
	odataConfig = config
}
