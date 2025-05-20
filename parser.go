package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Parser struct {
	Files []File
}

func NewParser() *Parser {
	return &Parser{
		Files: make([]File, 0),
	}
}

type Project struct {
	Name     string
	Dirpath  string
	Packages []Package
}

type Package struct {
	Name  string
	Files []File
}

type File2 struct {
	Dirpath string
	Name    string
	Imports []string
}

type File struct {
	ParentDir string
	Package   string
	FileName  string
	Imports   []string
}

func (p *Parser) ParseFile(filePath string) error {
	var f File
	f.ParentDir = filepath.Dir(filePath)
	f.FileName = filepath.Base(filePath)
	f.Imports = make([]string, 0)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.SkipObjectResolution)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}
	f.Package = file.Name.Name
	for _, imp := range file.Imports {
		f.Imports = append(f.Imports, imp.Path.Value)
	}
	p.Files = append(p.Files, f)
	log.Printf("Parsed file: %+v", f)
	return nil
}

func (p *Parser) ParseDir(dirPath string) error {
	log.Println("Parsing directory:", dirPath)
	dirs, err := getAllDirs(dirPath)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		for _, pkg := range pkgs {
			for path, file := range pkg.Files {
				var f File
				f.ParentDir = filepath.Dir(path)
				f.FileName = filepath.Base(path)
				f.Package = file.Name.Name
				for _, imp := range file.Imports {
					f.Imports = append(f.Imports, imp.Path.Value)
				}
				p.Files = append(p.Files, f)
			}
		}
	}
	log.Printf("Parsed files: %+v", p.Files)
	return nil
}
func getAllDirs(root string) ([]string, error) {
	allDirs := make([]string, 0)
	fn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			allDirs = append(allDirs, path)
		}
		return nil
	}
	err := filepath.WalkDir(root, fn)

	if err != nil {
		return nil, err
	}
	return allDirs, nil
}
