package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"went/app/service/github"
	"went/registry"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"
)

var _ command.Interface = (*Init)(nil)

type Init struct {
	command.Base
}

func init() {
	registry.Command.Register(&Init{})
}

func (c *Init) GetHeader() command.Header {
	return command.Header{
		Use:   "init",
		Short: "Initialize a project",
		Long:  "Creates the base scaffolding for a Went project.",
	}
}

func (c *Init) Invoke() any {
	return c.run
}

var validProjectName = func(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	validName := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	if !validName.MatchString(s) {
		return fmt.Errorf("name must be lowercase letters, numbers, and underscores only")
	}
	return nil
}

var notEmpty = func(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("cannot be empty")
	}
	return nil
}

func (c *Init) run(
	console *clio.Clio,
	github github.ClientInterface,
) error {
	console.Banner()

	projectName := console.Ask("Project name", validProjectName)
	shortDesc := console.Ask("Short description", notEmpty)
	longDesc := console.Ask("Long description", notEmpty)

	template := console.MultipleChoice("Which kind?", []clio.Choice{
		{Value: "went-cli-template", Label: "CLI"},
		{Value: "went-web-template", Label: "Web"},
	})

	for {
		err := c.initProject(
			context.Background(),
			github,
			template,
			projectName,
			shortDesc,
			longDesc,
		)

		if err != nil {
			console.Error("Error initializing project: %s", err.Error())
			if console.Confirm("Retry?", true) {
				continue
			}

			return err
		}

		break
	}

	console.Success("Project '%s' initialized.", projectName)

	return nil
}

func (c *Init) initProject(
	ctx context.Context,
	github github.ClientInterface,
	template string,
	projectName string,
	projectShort string,
	projectLong string,
) error {
	releases, err := github.ListReleases(ctx, "zeroSal", template)
	if err != nil {
		return fmt.Errorf("failed to list releases: %w", err)
	}

	if len(releases) == 0 {
		return fmt.Errorf("no releases found")
	}

	latest := releases[0]
	if len(latest.Assets) == 0 {
		return fmt.Errorf("no assets found in release %s", latest.TagName)
	}

	data, err := github.DownloadAssetBytes(ctx, latest.Assets[0].BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	tr := tar.NewReader(gr)

	files := make(map[string][]byte)

	for {
		header, err := tr.Next()
		if err != nil {
			break
		}

		name := strings.TrimPrefix(header.Name, "template/")
		name = strings.ReplaceAll(name, "{{ PROJECT_NAME }}", projectName)

		switch header.Typeflag {
		case tar.TypeDir:
			if name == "" {
				continue
			}
			if err := os.MkdirAll(name, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", name, err)
			}
		case tar.TypeReg:
			fileData := make([]byte, header.Size)
			n, err := io.ReadFull(tr, fileData)
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read file %s: %w", name, err)
			}
			files[name] = fileData[:n]
		}
	}

	for filename, data := range files {
		replaced := string(data)
		replaced = strings.ReplaceAll(replaced, "{{ PROJECT_NAME }}", projectName)
		replaced = strings.ReplaceAll(replaced, "{{ SHORT_PROJECT_DESCRIPTION }}", projectShort)
		replaced = strings.ReplaceAll(replaced, "{{ LONG_PROJECT_DESCRIPTION }}", projectLong)

		if err := os.WriteFile(filename, []byte(replaced), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	return nil
}
