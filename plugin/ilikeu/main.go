package ilikeu

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/FloatTech/floatbox/file"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"strings"
)

var (
	engine = control.AutoRegister(&ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Brief:            "我喜欢你",
		Help:             "- 我喜欢你",
		PublicDataFolder: "ILikeU",
	})
	datapath = file.BOTPATH + "/" + engine.DataFolder()
)

func init() {
	rand.NewSource(time.Now().UnixNano())
	
	// 注册关键词响应
	engine.OnKeywordGroup([]string{"我喜欢你"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			image := randomImage()
			if image == "" {
				ctx.SendChain(message.Text("我喜欢你"))
				return
			}
			ctx.SendChain(message.Text("我喜欢你"), message.Image(image))
		})
}

func downloadAsset() {
	src := "./asset"
	err := copyDir(src, datapath)
	if err != nil {
		fmt.Printf("复制资源失败: %v\n", err)
	}
}

func copyDir(src, dst string) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(src, path)
			dstPath := filepath.Join(dst, rel)
			return copyFile(path, dstPath)
		}
		return nil
	})
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func randomImage() string {
	files, err := os.ReadDir(datapath)
	if err != nil || len(files) == 0 {
		return ""
	}

	var images []string
	for _, f := range files {
		if isImage(f.Name()) {
			images = append(images, f.Name())
		}
	}

	if len(images) == 0 {
		return ""
	}

	return "file:///" + filepath.Join(datapath, images[rand.Intn(len(images))])
}

func isImage(filename string) bool {
	ext := filepath.Ext(strings.ToLower(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return true
	default:
		return false
	}
}
