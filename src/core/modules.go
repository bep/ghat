package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/rs/zerolog/log"
)

func (myFlags *Flags) UpdateModule(file string) error {
	parser := hclparse.NewParser()
	src, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s", file)
	}

	parsedFile, fileDiags := parser.ParseHCL(src, file)

	if fileDiags.HasErrors() {
		return fmt.Errorf("config parse: %w", fileDiags)
	}

	content := parsedFile.Body.(*hclsyntax.Body)

	//newContent:=hcl.File{
	//	Body:  nil,
	//	Bytes: nil,
	//	Nav:   nil,
	//}

	newBody := hclsyntax.Body{}

	for _, block := range content.Blocks {

		if block.Type == "module" {
			var source *hclsyntax.Attribute
			var version *hclsyntax.Attribute
			//version := cty.Value{}
			//source := cty.Value{}

			myAttributes := block.Body.Attributes
			newAttributes := hclsyntax.Attributes{}

			if myAttributes["source"] != nil {
				source = myAttributes["source"]
			}

			if myAttributes["version"] != nil {
				version = myAttributes["version"]
			}

			log.Print(source)
			log.Print(version)
			for x, attribute := range myAttributes {

				if attribute.Name == "version" {
					continue
				}

				if attribute.Name == "source" {
					ctx := &hcl.EvalContext{}

					var diags hcl.Diagnostics
					sourceValue, diags := attribute.Expr.Value(ctx)

					myFlags.GetType(sourceValue.AsString())
					// what is the source type
					//
					// if source type is git  update version?
					// if source is registry convert to git
					if diags.HasErrors() {
						return fmt.Errorf("version parse: %w", fileDiags)
					}
				}

				newAttributes[x] = attribute
				//		ctx := &hcl.EvalContext{}
				//		var diags hcl.Diagnostics
				//		version, diags = attribute.Expr.Value(ctx)
				//
				//		if diags.HasErrors() {
				//			return nil, fmt.Errorf("version parse: %w", fileDiags)
				//		}
				//	}
			}

			block.Body.Attributes = newAttributes
			//log.Info().Msgf("%s", version)
			//new tf file

			//terraformBlock, err := p.parseBlock(block, file)
			//src, _ := os.ReadFile(file)
			//
			//hclFile, diags := hclwrite.ParseConfig(src, file, hcl.InitialPos)
			//
			//if diags.HasErrors() {
			//	return nil, fmt.Errorf("config parse: %w", diags)
			//}
			//
			//hclSyntaxFile, diags := hclsyntax.ParseConfig(src, file, hcl.InitialPos)

			//f.Body().SetAttributeValue("version", cty.StringVal("999.999.999"))

			//myBlocks := test.Blocks()
			//for _, block := range myBlocks {
			//	log.Print(block.Labels())
			//}
			//log.Print(myBlocks)
			//result := hclFile.Bytes()
			//log.Print(result)
			//log.Print(hclSyntaxFile)
			newBody.Blocks = append(newBody.Blocks, block)
		} else {
			newBody.Blocks = append(newBody.Blocks, block)
			log.Print(newBody)
		}

	}

	//write back if modified and if !myFlags.DryRun
	return nil
}

func (myFlags *Flags) UpdateModules() error {

	terraform, err := myFlags.GetTF()

	if err != nil {
		return err
	}

	// contains a module?
	for _, file := range terraform {
		err = myFlags.UpdateModule(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (myFlags *Flags) GetTF() ([]string, error) {
	var terraform []string

	for _, match := range myFlags.Entries {
		//for each file that is a terraform file
		if path.Ext(match) == ".tf" {
			terraform = append(terraform, match)
		}
	}

	return terraform, nil
}

func (myFlags *Flags) GetType(module string) (string, string, error) {
	var moduleType string

	//handle local path
	absPath, _ := filepath.Abs(module)
	_, err := os.Stat(absPath)

	if err == nil {
		return module, "local", nil
	}

	if strings.Contains("?ref=", module) {
		moduleType = "url"

		if myFlags.Update {
			splitter := strings.Split(module, "?ref=")
			base := splitter[0]
			log.Print(base)
			// get lastest tag from git reference
			//trim git:
			//trim
		}

		return module, moduleType, nil
	}

	//split and check if using version then replace with hash or updated version hash
	//then it's a good
	//Modules in Package Sub-directories
	//
	//Generic Git, Mercurial repositories
	//
	//	//HTTP URLs
	//	//	//
	//	//GitHub

	//Local paths
	//
	//Terraform Registry

	//
	//Bitbucket

	//S3 buckets
	//
	//GCS buckets
	//

	return module, moduleType, err
}
