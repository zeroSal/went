package github

import (
	"context"
	"io"
	"went/app/model/github"
)

type ClientInterface interface {
	ListReleases(
		ctx context.Context,
		owner,
		repo string,
	) ([]github.Release, error)

	DownloadAsset(
		ctx context.Context,
		downloadURL string,
		writer io.Writer,
	) error

	DownloadAssetBytes(
		ctx context.Context,
		downloadURL string,
	) ([]byte, error)
}
