package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apstndb/spantype"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"strings"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"log"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func modeToFormatOption(mode string) spantype.FormatOption {
	switch strings.ToLower(mode) {
	case "verbose":
		return spantype.FormatOptionVerbose
	case "normal":
		return spantype.FormatOptionNormal
	case "simplest":
		return spantype.FormatOptionSimplest
	case "simple":
		return spantype.FormatOptionSimple
	default:
		panic("unknown mode: " + mode)
	}
}

func run(ctx context.Context) error {
	mode := flag.String("mode", "verbose", "format mode (simplest|simple|normal|verbose)")
	flag.Parse()

	formatOpt := modeToFormatOption(*mode)

	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var structType sppb.StructType
	if err := protojson.Unmarshal(b, &structType); err != nil {
		return err
	}
	fmt.Println(spantype.FormatStructFields(structType.GetFields(), formatOpt))
	return nil
}
