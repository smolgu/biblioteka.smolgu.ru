package vk

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/smolgu/lib/modules/setting"
	"github.com/zhuharev/vkutil"
)

var (
	api = vkutil.New()
)

func getPhotoString(photoID, ownerID int) string {
	return fmt.Sprintf("photo%d_%d", ownerID, photoID)
}

// Post send post to vk wall
// returns post id
func Post(message string, photos ...vkutil.Photo) (int, error) {
	api.VkApi.AccessToken = setting.VkAccessToken

	api.SetDebug(true)

	params := url.Values{}
	if len(photos) > 0 {
		att := []string{}
		for _, v := range photos {
			att = append(att, getPhotoString(v.Id, v.OwnerId))
		}
		params.Set("attachments", strings.Join(att, ","))
	}

	return api.WallPost(vkutil.OptsWallPost{
		OwnerId:   -setting.VkGroupID,
		Message:   message,
		FromGroup: true,
	}, params)
}

// Upload upload image to vk server, return id
func Upload(path string) (vkutil.Photo, error) {
	api.SetDebug(true)
	api.VkApi.AccessToken = setting.VkAccessToken

	f, err := os.Open(filepath.Join("./data/", path))
	if err != nil {
		return vkutil.Photo{}, err
	}

	return api.PhotosUploadWall(setting.VkGroupID, f)
}
