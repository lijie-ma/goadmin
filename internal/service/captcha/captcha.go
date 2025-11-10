package captcha

import (
	"encoding/json"
	"fmt"
	"goadmin/internal/context"
	"goadmin/internal/model/captcha"
	"goadmin/internal/service/setting"
	"goadmin/pkg/redisx"
	"goadmin/pkg/util"
	"sync"
	"time"

	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

var (
	sliderMutex    sync.Mutex
	slideBasicCapt slide.Captcha
	logPrefix      = "captcha"
)

func getCapt() (slide.Captcha, error) {
	sliderMutex.Lock()
	defer sliderMutex.Unlock()

	if slideBasicCapt != nil {
		return slideBasicCapt, nil
	}

	builder := slide.NewBuilder(
	//slide.WithGenGraphNumber(2),
	//slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	imgs, err := imagesv2.GetImages()
	if err != nil {
		return nil, err
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		return nil, err
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	slideBasicCapt = builder.Make()
	return slideBasicCapt, nil
}

func Generate(ctx *context.Context) (any, error) {
	captchaCfg, err := setting.GetCaptchaSwitch(ctx, setting.NewServerSettingService())
	if err != nil {
		ctx.Logger.Errorf("%s GetCaptchaSwitch %+v", logPrefix, err)
		return nil, err
	}
	if !captchaCfg.IsAdminOn() {
		return map[string]any{
			"switch": 0,
		}, nil
	}

	capt, err := getCapt()
	if err != nil {
		ctx.Logger.Errorf("%s capt %+v", logPrefix, err)
		return nil, err
	}
	captData, err := capt.Generate()
	if err != nil {
		ctx.Logger.Errorf("%s Generate %+v", logPrefix, err)
		return nil, err
	}
	blockData := captData.GetData()
	if blockData == nil {
		ctx.Logger.Errorf("%s GetData %+v", logPrefix, err)
		return nil, err
	}
	var masterImageBase64, tileImageBase64 string
	masterImageBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		ctx.Logger.Errorf("%s GetMasterImage %+v", logPrefix, err)
		return nil, err
	}
	tileImageBase64, err = captData.GetTileImage().ToBase64()
	if err != nil {
		ctx.Logger.Errorf("%s GetTileImage %+v", logPrefix, err)
		return nil, err
	}
	dotsByte, _ := json.Marshal(blockData)
	key := util.GenerateUUIDWithoutHyphen()
	redisx.GetClient().Set(ctx, key, string(dotsByte), time.Minute)

	return map[string]any{
		"switch":       1,
		"key":          key,
		"image_base64": masterImageBase64,
		"tile_base64":  tileImageBase64,
		"tile_width":   blockData.Width,
		"tile_height":  blockData.Height,
		"tile_x":       blockData.DX,
		"tile_y":       blockData.DY,
	}, nil
}

func Check(ctx *context.Context, formData captcha.CheckForm) error {
	redisClient := redisx.GetClient()
	catchData, err := redisClient.Get(ctx, formData.Key).Result()
	if err != nil {
		ctx.Logger.Errorf("%s redis Get %+v", logPrefix, err)
		return err
	}
	redisClient.Del(ctx, formData.Key)

	var dct *slide.Block
	if err := json.Unmarshal([]byte(catchData), &dct); err != nil {
		ctx.Logger.Errorf("%s Unmarshal catchData %s %+v", logPrefix, catchData, err)
		return err
	}

	if !slide.Validate(formData.X, formData.Y, dct.X, dct.Y, 5) {
		ctx.Logger.Warnf("%s Validate failed %+v %+v", logPrefix, formData, dct)
		return fmt.Errorf("%s Validate failed", logPrefix)
	}
	return nil
}
