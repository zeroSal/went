package cmd

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"went/app"
	"went/app/service/template"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type InitCmd struct {
	embedFS    embed.FS
	buildSpecs *app.BuildSpecs
}

func NewInitCmd(
	embedFS embed.FS,
	buildSpecs *app.BuildSpecs,
) *InitCmd {
	return &InitCmd{
		embedFS:    embedFS,
		buildSpecs: buildSpecs,
	}
}

func (s *InitCmd) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Run:   s.run,
	}
}

func (s *InitCmd) run(cmd *cobra.Command, args []string) {
	a := fx.New(
		fx.Supply(s.embedFS),
		fx.Provide(func() *app.BuildSpecs {
			return s.buildSpecs
		}),
		app.Kernel,
		fx.Invoke(s.init),
	)

	ctx := context.Background()
	if err := a.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start: %v\n", err)
		return
	}
	a.Stop(ctx)
}

func (s *InitCmd) init() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	defaultName := filepath.Base(cwd)
	defaultName = strings.ReplaceAll(defaultName, "-", "_")

	fmt.Printf("Project name (default: %s): ", defaultName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	projectName := strings.TrimSpace(input)
	if projectName == "" {
		projectName = defaultName
	}

	fmt.Printf("Initializing project: %s\n", projectName)

	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod init: %w", err)
	}

	renderer := template.NewRenderer(s.embedFS)

	data := map[string]string{
		"ProjectName": projectName,
	}

	if err := renderer.Render(data); err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}

	for _, file := range renderer.GetFiles() {
		content, ok := renderer.Read(file)
		if !ok {
			continue
		}

		outputPath := filepath.Join(cwd, file)
		dir := filepath.Dir(outputPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(outputPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}

		fmt.Printf("Created: %s\n", outputPath)
	}

	fmt.Println("\nProject initialized successfully!")
	fmt.Println("\nTo build and run:")
	fmt.Println("  go mod tidy")
	fmt.Printf("  go build -o %s .\n", projectName)
	fmt.Printf("  ./%s serve\n", projectName)
	return nil
}
