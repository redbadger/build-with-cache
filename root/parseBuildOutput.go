package root

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

func parseBuildOutput(buildOutput, imgTag string, stages []string) (refs map[string]string, err error) {
	var ref reference.Named
	prev := ""
	stageNum := 0
	refs = make(map[string]string)

	for _, s := range strings.Split(buildOutput, "\n") {
		if strings.Contains(s, ": FROM ") || strings.Contains(s, "Successfully built") {
			if strings.HasPrefix(prev, " ---> ") {
				sha := strings.Fields(prev)[1]
				ref, err = reference.ParseNamed(imgTag)
				if err != nil {
					return
				}
				stage := stages[stageNum]
				ref := fmt.Sprintf("%s-%s", reference.TrimNamed(ref), stage)
				stageNum++
				refs[sha] = ref
			}
		}
		prev = s
	}
	return
}
