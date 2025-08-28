package thumb

import (
	"fmt"
	"lyworder/pkg/setting"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// ThumbnailGenerator 缩略图生成器
type ThumbnailGenerator struct {
	config *setting.ThumbSettingS
}

// NewThumbnailGenerator 创建新的缩略图生成器
func NewThumbnailGenerator(config *setting.ThumbSettingS) *ThumbnailGenerator {
	return &ThumbnailGenerator{
		config: config,
	}
}

// GenerateThumbnail 异步生成缩略图
func (tg *ThumbnailGenerator) GenerateThumbnail(imagePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}

	// 提取文件名和扩展名
	fileName := filepath.Base(imagePath)
	fileExt := filepath.Ext(fileName)
	fileNameWithoutExt := fileName[:len(fileName)-len(fileExt)]

	// 构建缩略图路径
	thumbDir := filepath.Dir(imagePath)
	thumbPath := filepath.Join(thumbDir, fmt.Sprintf("%s.thumb%s", fileNameWithoutExt, fileExt))

	// 异步处理
	go func() {
		// 读取图片
		img, err := imaging.Open(imagePath)
		if err != nil {
			fmt.Printf("failed to open image: %v\n", err)
			return
		}

		// 调整图片大小
		thumb := imaging.Fit(img, tg.config.MaxWidth, tg.config.MaxHeight, imaging.Lanczos)

		// 保存缩略图
		if err := imaging.Save(thumb, thumbPath); err != nil {
			fmt.Printf("failed to save thumbnail: %v\n", err)
			return
		}

		fmt.Printf("thumbnail generated: %s\n", thumbPath)
	}()

	return nil
}

// 检查图片类型
func isImageFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
}
