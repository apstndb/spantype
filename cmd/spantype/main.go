package main

import (
	"context"
	"fmt"
	"github.com/apstndb/spantype"
	"google.golang.org/protobuf/encoding/protojson"
	"io"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"log"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var structType sppb.StructType
	if err := protojson.Unmarshal(b, &structType); err != nil {
		return err
	}
	fmt.Println(spantype.FormatStructFields(structType.GetFields(), spantype.FormatTypeOptionVerbose))
	return nil
}
