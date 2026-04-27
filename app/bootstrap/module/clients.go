package module

import (
	"context"

	"went/app/service/github"
	"went/app/service/http"

	"go.uber.org/fx"
)

func HTTPClientProvider() http.ClientInterface {
	return http.NewClient(nil)
}

func GitHubClientProvider(
	lc fx.Lifecycle,
	httpClient http.ClientInterface,
) github.ClientInterface {
	client := github.NewClient(httpClient)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return client
}
