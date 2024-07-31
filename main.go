package main

import (
	"context"
	"fmt"
	"github.com/kr/pretty"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	cwd, _ := os.Getwd()
	ctx := context.Background()
	err := os.Setenv("PULUMI_CONFIG_PASSPHRASE", "PULUMI_CONFIG_PASSPHRASE")
	if err != nil {
		fmt.Printf("Error setting environment variable: %v\n", err)
		return
	}

	projectWorkspace := auto.Project(workspace.Project{
		Name:    "example",
		Runtime: workspace.NewProjectRuntimeInfo("go", nil),
		Backend: &workspace.ProjectBackend{
			URL: "file://" + cwd,
		},
	})

	stack, err := auto.UpsertStackInlineSource(ctx, "qa", "example", runFunc, projectWorkspace)
	if err != nil {
		fmt.Printf("Cannot create stack: %v\n", err)
		return
	}

	ws := stack.Workspace()
	err = ws.InstallPlugin(ctx, "random", "v4.16.3")
	if err != nil {
		fmt.Printf("Error installing plugin: %v\n", err)
		return
	}

	_, err = stack.Refresh(ctx)
	if err != nil {
		fmt.Printf("Cannot refresh stack: %v\n", err)
		return
	}

	result, err := stack.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		fmt.Printf("Stack up error: %v\n", err)
		return
	}

	pretty.Println(result)
}

func runFunc(ctx *pulumi.Context) error {
	password, err := random.NewRandomPassword(ctx, "password", &random.RandomPasswordArgs{
		Length:          pulumi.Int(16),
		Special:         pulumi.Bool(true),
		OverrideSpecial: pulumi.String("!#$%&*()-_=+[]{}<>:?"),
	})
	if err != nil {
		return err
	}
	ctx.Export("password", password.Result)
	return nil
}
