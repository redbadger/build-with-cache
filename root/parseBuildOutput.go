package root

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

func parseBuildOutput(buildOutput, imgName, cache string, stages []string) (imgNames map[string]string, err error) {
	prev := ""
	stageNum := 0
	imgNames = make(map[string]string)

	imgNamed, err := getNamed(imgName, cache)
	if err != nil {
		return
	}
	for _, s := range strings.Split(buildOutput, "\n") {
		if strings.Contains(s, ": FROM ") || strings.Contains(s, "Successfully built") {
			if strings.HasPrefix(prev, " ---> ") {
				sha := strings.Fields(prev)[1]
				stage := stages[stageNum]
				imgNamed := fmt.Sprintf("%s-%s", reference.TrimNamed(imgNamed), stage)
				stageNum++
				imgNames[sha] = imgNamed
			}
		}
		prev = s
	}
	return
}
