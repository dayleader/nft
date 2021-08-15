package main

import (
	"flag"
	"fmt"
	"nft/internal/combiner"
	"nft/internal/domain"
	"nft/internal/generator"
	"nft/internal/inmemory"
	"nft/internal/trait"
	"os"
)

func main() {
	// Flags.
	appContext := domain.NewAppContext()
	flag.IntVar(&appContext.GeneratorParams.Length, "generated-image-length", 2048, "generated image length")
	flag.IntVar(&appContext.GeneratorParams.Width, "generated-image-width", 2048, "generated image width")
	flag.StringVar(&appContext.GeneratorParams.InputDirectory, "generated-image-input", "input-dir", "generated image input directory")
	flag.StringVar(&appContext.GeneratorParams.OutputDirectory, "generated-image-output", "output-dir", "generated image output directory")
	flag.IntVar(&appContext.GeneratorParams.Number, "generated-image-number", 100, "generated image number")
	generate := flag.Bool("generate", false, "Generate images")

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	if *generate {
		fmt.Println("Generating images ...")
		traitService := trait.NewBasicTraitService(
			inmemory.NewGroupRepository(),
			inmemory.NewTraitRepository(),
		)
		_, err := traitService.Import(appContext.GeneratorParams.InputDirectory)
		if err != nil {
			panic(err)
		}
		generatorService := generator.NewBasicImageGenerator(
			appContext.GeneratorParams,
			traitService,
			combiner.NewBasicImageCombiner(),
		)
		err = generatorService.GenerateImages()
		if err != nil {
			panic(err)
		}
		return
	}
}
