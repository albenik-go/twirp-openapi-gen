package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/albenik/twirp-openapi-gen/internal/openapi20"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	fs := new(flag.FlagSet)
	hostname := fs.String("hostname", "example.com", "")
	pathPrefix := fs.String("path_prefix", "", "")
	outputSuffix := fs.String("output_suffix", ".swagger.json", "")

	protogen.Options{ParamFunc: fs.Set}.
		Run(func(gen *protogen.Plugin) error {
			gen.SupportedFeatures = 1 // proto3 optional fields support
			openapiGen := openapi20.NewGenerator()

			for _, file := range gen.Files {
				if !file.Generate {
					continue
				}

				log.Println("processing:", file.Desc.Path())

				for _, service := range file.Services {
					schema, err := openapiGen.GenerateSchema(*hostname, *pathPrefix, service)
					if err != nil {
						return fmt.Errorf("%s: schema: %w", file.Desc.Path(), err)
					}

					fname := filepath.Join(filepath.Dir(file.GeneratedFilenamePrefix),
						string(service.Desc.Name())+*outputSuffix)
					j := json.NewEncoder(gen.NewGeneratedFile(fname, file.GoImportPath))
					j.SetIndent("", "  ")
					if err = j.Encode(schema); err != nil {
						return fmt.Errorf("%s: encode: %w", file.Desc.Path(), err)
					}
				}
			}

			return nil
		})
}
