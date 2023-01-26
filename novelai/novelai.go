package novelai

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-json"
	"github.com/urfave/cli/v2"
)

// NOTE: NovelAI inside prompt
// const defaultPrompt = "masterpiece, best quality, "
// const noneUndesiredContext = "lowres"
// const lowQualityUndesiredContext = "nsfw, lowres, text, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry"
// const lowQualityAndBadAnatomyUndesiredContext = "nsfw, lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry"

const aiGeneratedSignature = "AI generated image"
const descriptionStart = "Description"
const descriptionEnd = string(uint8(00)) + string(uint8(00)) + string(uint8(00)) + string(uint8(16))
const commentStart = "EXtComment"
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

func CheckDirectory(cCtx *cli.Context) error {
	dirName := cCtx.Args().Get(0)
	err := filepath.WalkDir(dirName, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		fmt.Println(info.Name())
		err = checkFile(path)
		if err != nil {
			return err
		}
		fmt.Println("")

		return nil
	})
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}

func checkFile(fileName string) error {
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

	// img size
	reader, err := os.Open(fileName)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	defer reader.Close()
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	fmt.Println("Size:", img.Width, "x", img.Height)

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
