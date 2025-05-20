package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	fileFlag = "file"
	dirFlag  = "dir"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	app := &cli.App{
		Name:  "VisAI",
		Usage: "Code visualization and analysis tool",
		Commands: []*cli.Command{
			{
				Name:    "parse",
				Aliases: []string{"p"},
				Usage:   "Parse the codebase and store relationships",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     dirFlag,
						Aliases:  []string{"d"},
						Usage:    "Project root folder to begin parsing",
						Required: false,
					},
					&cli.StringFlag{
						Name:     fileFlag,
						Aliases:  []string{"f"},
						Usage:    "Individual File to parse",
						Required: false,
					},
				},
				Before: validateFlags,
				Action: parseAction,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func validateFlags(c *cli.Context) error {
	if c.IsSet("dir") && c.IsSet("file") {
		return fmt.Errorf("cannot use both --dir and --file flags together")
	}
	if !c.IsSet("dir") && !c.IsSet("file") {
		return fmt.Errorf("either --dir or --file must be specified")
	}
	if c.IsSet("dir") {
		dirPath := c.String("dir")
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		isDir, err := validateIsDir(absPath)
		if err != nil {
			return fmt.Errorf("failed to validate directory: %w", err)
		}
		if !isDir {
			return fmt.Errorf("provided path %s is not a directory", absPath)
		}
	}
	if c.IsSet("file") {
		filePath := c.String("file")
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		isFile, err := validateIsFile(absPath)
		if err != nil {
			return fmt.Errorf("failed to validate file: %w", err)
		}
		if !isFile {
			return fmt.Errorf("provided path %s is not a file", absPath)
		}
	}
	return nil
}
func validateIsFile(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Mode().IsRegular(), nil
}
func validateIsDir(dirPath string) (bool, error) {

	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.IsDir(), nil
}
func parseAction(c *cli.Context) error {
	parser := NewParser()
	if c.IsSet(fileFlag) {
		filePath := c.String(fileFlag)
		absoluteFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		err = parser.ParseFile(absoluteFilePath)
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		return nil
	}
	if c.IsSet(dirFlag) {
		projectRoot := c.String(dirFlag)
		absoluteRootPath, err := filepath.Abs(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
		log.Println("Project root:", absoluteRootPath)
		if err := parser.ParseDir(absoluteRootPath); err != nil {
			return fmt.Errorf("failed to parse directory: %w", err)
		}
		return nil
	}
	return nil
}
