package provider

import "terraform-provider-coolify/internal/api"

type CoolifyProviderData struct {
	endpoint string
	client   *api.APIClient
}
