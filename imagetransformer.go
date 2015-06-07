package main

import "os"
import "fmt"
import "flag"
import "path/filepath"
import "image"
import (
	"image/jpeg"
	_ "image/jpeg"
)
import "github.com/disintegration/imaging"

var sourceDir string
var targetDir string
var transformation string
var param int

func init() {
	flag.StringVar(&sourceDir, "sourceDir", ".", "Source directory where to look for file to convert")
	flag.StringVar(&targetDir, "targetDir", ".", "Target directory where to put the converted files")
	flag.StringVar(&transformation, "transformation", "resize", "Transformation to apply: resize (default) or crop")
	flag.IntVar(&param, "param", 3, "")
}

func main() {
	flag.Parse()
	flag.Visit(func(flag* flag.Flag) {
		fmt.Println(flag.Name)
	})
	var t Transformer
	switch transformation {
	case "resize":
		t = ResizeTransformer{param}
	case "crop":
		t = CropTransformer{param}
	default:
		fmt.Println("Unsupported transformation:", transformation)
	}
	w := Walker{t}
	filepath.Walk(sourceDir, w.WalkFunc)
}

type Walker struct {
	Former Transformer
}

func (w *Walker) WalkFunc(path string, info os.FileInfo, err error) error {
	// short-circuit for directories
	fmt.Println("Now dealing with:", path)
	if fid, fiderr := os.Stat(path); fiderr == nil {
		if fid.IsDir() {
			fmt.Println(path, "is a directory, skipping")
			return nil
		}
	}

	// Magical manipulation to get the new path for the converted image
	imgRelPath, _ := filepath.Rel(sourceDir, path)
	targetImgPath := filepath.Join(targetDir, imgRelPath)
	targetImgDir := filepath.Dir(targetImgPath)
	fmt.Println("Trying to create:", targetImgPath)
	if _, existErr := os.Stat(targetImgPath); existErr == nil {
		fmt.Println("Already present:", targetImgPath)
	} else {
		if mkErr := os.MkdirAll(targetImgDir, os.ModeDir); mkErr != nil {
			fmt.Println("Unable to create:", targetImgDir, mkErr)
			return mkErr
		}
		return w.Former.Transform(path, targetImgPath)
	}

	return nil
}

type Transformer interface {
	Transform(sourceFilename, targetFilename string) error
}

type ResizeTransformer struct {
	Factor int
}

func (t ResizeTransformer) Transform(sourceFilename, targetFilename string) error {
	fmt.Println("Converting:", sourceFilename, targetFilename)
	reader, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer reader.Close()
	img, _, _ := image.Decode(reader)
	v := t.view(img, t.Factor)
	thumb := imaging.Resize(img, v.Dx(), v.Dy(), imaging.Lanczos)
	writer, err := os.Create(targetFilename)
	if err != nil {
		return err
	}
	err = jpeg.Encode(writer, thumb, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t ResizeTransformer) view(img image.Image, cs int) image.Rectangle {
	ir := img.Bounds()
	return image.Rect(0, 0, ir.Dx()/cs, ir.Dy()/cs)
}

type CropTransformer struct {
	Size int
}

func (t CropTransformer) Transform(sourceFilename, targetFilename string) error {
	fmt.Println("Converting:", sourceFilename, targetFilename)
	reader, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer reader.Close()
	img, _, _ := image.Decode(reader)
	v := t.view(img, t.Size)
	thumb := imaging.Crop(img, v)
	writer, err := os.Create(targetFilename)
	if err != nil {
		return err
	}
	err = jpeg.Encode(writer, thumb, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t CropTransformer) view(img image.Image, cs int) image.Rectangle {
	ir := img.Bounds()
	csx := min(ir.Dx(), cs)
	csy := min(ir.Dy(), cs)

	lx := ir.Min.X + (ir.Dx()-csx)/2
	ly := ir.Min.Y + (ir.Dy()-csy)/2

	return image.Rect(lx, ly, lx+csx, ly+csy)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
