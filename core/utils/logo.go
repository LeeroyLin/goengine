package utils

import (
	"strings"
)

func GetLogoStr() string {
	builder := strings.Builder{}

	spacing := "          "

	builder.WriteString("\n")
	builder.WriteString(spacing)
	builder.WriteString("▉▉▉▉    ▉▉▉       ▉▉▉    ▉  ▉    ▉▉▉    ▉   ▉    ▉▉▉▉    ▉▉▉    ▉▉▉▉\n")
	builder.WriteString(spacing)
	builder.WriteString("▉  ▉    ▉  ▉     ▉       ▉  ▉     ▉     ▉▉ ▉▉    ▉        ▉     ▉  ▉\n")
	builder.WriteString(spacing)
	builder.WriteString("▉▉▉▉    ▉▉▉      ▉       ▉▉▉▉     ▉     ▉ ▉ ▉    ▉▉▉▉     ▉     ▉▉▉▉\n")
	builder.WriteString(spacing)
	builder.WriteString("▉  ▉    ▉  ▉     ▉       ▉  ▉     ▉     ▉   ▉    ▉        ▉     ▉  ▉\n")
	builder.WriteString(spacing)
	builder.WriteString("▉  ▉    ▉   ▉     ▉▉▉    ▉  ▉    ▉▉▉    ▉   ▉    ▉▉▉▉     ▉     ▉  ▉\n")

	return builder.String()
}
