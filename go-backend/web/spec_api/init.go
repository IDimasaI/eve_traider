// file: specapi/specapi.go
package specapi

import (
	"local_server/web/spec_api/cloud"
	"local_server/web/spec_api/desktop"
)

type Platform struct {
	Cloud   *cloud.Cloud
	Desktop *desktop.Desktop
}

func CreatePlatform(useCloud bool) *Platform {
	if useCloud {
		c := cloud.NewCloud()

		return &Platform{
			Cloud: c,
		}
	} else {
		d := desktop.NewDesktop()

		return &Platform{
			Desktop: d,
		}
	}
}
