package novelai

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/urfave/cli/v2"
)

const aiGeneratedSignature = "AI generated image"
const descriptionStart = "tEXtDescription"
const descriptionEnd = string(uint8(00)) + string(uint8(00)) + string(uint8(00)) + string(uint8(16))
const commentStart = "tEXtComment"
const commentEnd = string(uint8(00)) + string(uint8(01)) + string(uint8(00)) + string(uint8(00))

type Result struct {
	Prompt  string
	Comment Comment
}

type Comment struct {
	Steps            int64   `json:"steps"`
	Sampler          string  `json:"sampler"`
	Seed             int64   `json:"seed"`
	Strength         float64 `json:"strength"`
	Noise            float64 `json:"noise"`
	Scale            float64 `json:"scale"`
	UndesiredContent string  `json:"uc"`
}

func CheckFile(cCtx *cli.Context) error {
	fileName := cCtx.Args().Get(0)
	f, err := os.ReadFile(fileName)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	result, err := getResult(f)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// print
	fmt.Println("Prompt:", result.Prompt)
	fmt.Println("Undesired Content:", result.Comment.UndesiredContent)
	fmt.Println("Steps:", result.Comment.Steps)
	fmt.Println("Scale:", result.Comment.Scale)
	fmt.Println("Seed:", result.Comment.Seed)
	fmt.Println("Sampling:", result.Comment.Sampler)
	fmt.Println("(hidden) Strength:", result.Comment.Strength)
	fmt.Println("(hidden) Noise:", result.Comment.Noise)

	return nil
}

func getResult(f []byte) (*Result, error) {
	var result = &Result{}
	str := string(f)
	isNovelAIImage := strings.Contains(str, aiGeneratedSignature)
	if isNovelAIImage {
		// get description
		descStartIdx := strings.Index(str, descriptionStart) + len(descriptionStart) + 1
		descEndIdx := strings.Index(str, descriptionEnd)
		prompt := string(f[descStartIdx : descEndIdx-len(descriptionEnd)])
		// get comment
		cmtStartIdx := strings.Index(str, commentStart) + len(commentStart) + 1
		cmtEndIdx := strings.Index(str, commentEnd)
		jsonData := string(f[cmtStartIdx : cmtEndIdx-len(commentEnd)])
		// set
		var comment Comment
		err := json.Unmarshal([]byte(jsonData), &comment)
		if err != nil {
			return nil, err
		}
		result.Prompt = prompt
		result.Comment = comment
	} else {
		return nil, errors.New("this image does not have an AI generated image signature")
	}
	return result, nil
}
