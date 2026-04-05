package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/apstndb/spantype"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func modeToFormatOption(mode string) (spantype.FormatOption, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "more":
		return spantype.FormatOptionMoreVerbose, nil
	case "verbose":
		return spantype.FormatOptionVerbose, nil
	case "normal":
		return spantype.FormatOptionNormal, nil
	case "simplest":
		return spantype.FormatOptionSimplest, nil
	case "simple":
		return spantype.FormatOptionSimple, nil
	default:
		return spantype.FormatOption{}, fmt.Errorf("unknown mode %q (want simplest|simple|normal|verbose|more)", mode)
	}
}

func parseTypeAnnotationMode(s string) (spantype.TypeAnnotationMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "suffix", "":
		return spantype.TypeAnnotationModeSuffix, nil
	case "omit":
		return spantype.TypeAnnotationModeOmit, nil
	case "primary":
		return spantype.TypeAnnotationModePrimary, nil
	default:
		return 0, fmt.Errorf("unknown type-annotation mode %q (want suffix|omit|primary)", s)
	}
}

func run() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	mode := fs.String("mode", "verbose", "format mode (simplest|simple|normal|verbose|more)")
	typeAnn := fs.String("type-annotation", "suffix", "how to render TypeAnnotation: suffix|omit|primary")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	formatOpt, err := modeToFormatOption(*mode)
	if err != nil {
		return err
	}
	annMode, err := parseTypeAnnotationMode(*typeAnn)
	if err != nil {
		return err
	}
	formatOpt.TypeAnnotation = annMode

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
