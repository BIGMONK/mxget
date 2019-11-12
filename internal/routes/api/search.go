package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/provider/baidu"
	"github.com/winterssy/mxget/pkg/provider/kugou"
	"github.com/winterssy/mxget/pkg/provider/kuwo"
	"github.com/winterssy/mxget/pkg/provider/migu"
	"github.com/winterssy/mxget/pkg/provider/netease"
	"github.com/winterssy/mxget/pkg/provider/qq"
	"github.com/winterssy/mxget/pkg/provider/xiami"
)

func SearchSongsFromNetEase(c *gin.Context) {
	searchSongs(c, netease.Client())
}

func SearchSongsFromQQ(c *gin.Context) {
	searchSongs(c, qq.Client())
}

func SearchSongsFromMiGu(c *gin.Context) {
	searchSongs(c, migu.Client())
}

func SearchSongsFromKuGou(c *gin.Context) {
	searchSongs(c, kugou.Client())
}

func SearchSongsFromKuWo(c *gin.Context) {
	searchSongs(c, kuwo.Client())
}

func SearchSongsFromXiaMi(c *gin.Context) {
	searchSongs(c, xiami.Client())
}

func SearchSongsFromBaiDu(c *gin.Context) {
	searchSongs(c, baidu.Client())
}

func searchSongs(c *gin.Context, client provider.API) {
	keyword := strings.TrimSpace(c.Param("keyword"))
	data, err := client.SearchSongs(keyword)
	if err != nil {
		c.JSON(500, &provider.Response{
			Code:     500,
			Msg:      err.Error(),
			Platform: client.PlatformId(),
		})
		return
	}

	c.JSON(200, &provider.Response{
		Code:     200,
		Data:     data,
		Platform: client.PlatformId(),
	})
}
