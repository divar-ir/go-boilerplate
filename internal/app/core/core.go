package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.cafebazaar.ir/bardia/lazyapi/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"git.cafebazaar.ir/bardia/lazyapi/pkg/appdetail"
)

type core struct {
	globalLock sync.Mutex
}

func New() appdetail.AppDetailServer {
	return &core{
		globalLock: sync.Mutex{},
	}
}

func (c *core) GetAppDetail(ctx context.Context, request *appdetail.GetAppDetailRequest) (*appdetail.GetAppDetailReply, error) {
	c.globalLock.Lock()
	defer c.globalLock.Unlock()
	time.Sleep(1 * time.Second)

	resp, err := http.Post("https://api.cafebazaar.ir/rest-v1/process",
		"application/json",
		bytes.NewBufferString(fmt.Sprintf("{\"singleRequest\":{\"appDetailsRequest\":{\"packageName\":\"%s\"}}}", request.PackageName)))
	if err != nil {
		return nil, errors.Wrap(err, "fail to send request to cafebazaar")
	}

	if resp.StatusCode != 200 {
		return &appdetail.GetAppDetailReply{
			StatusCode: int32(resp.StatusCode),
		}, nil
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to send request to cafebazaar")
	}
	var objmap = make(map[string]interface{})

	err = json.Unmarshal(b, &objmap)
	if err != nil {
		return nil, errors.Wrap(err, "fail to send request to cafebazaar")
	}

	respMap := objmap["singleReply"].(map[string]interface{})["appDetailsReply"].(map[string]interface{})
	packageMap:= respMap["package"].(map[string]interface{})
	var permissions []string
	for _,a := range packageMap["permissions"].([]interface{}){
		permissions = append(permissions, a.(string))
	}
	return &appdetail.GetAppDetailReply{
		Detail: &appdetail.App{
			Name:         respMap["name"].(string),
			Description:  respMap["description"].(string),
			Homepage:     respMap["homepage"].(string),
			Email:        respMap["appEmail"].(string),
			AuthorName:   respMap["authorName"].(string),
			CategoryName: respMap["categoryName"].(string),
			Package: &appdetail.Package{
				PackageId: uint32(packageMap["packageID"].(float64)),
				PackageHash: packageMap["packageHash"].(string),
				VersionCode: uint32(packageMap["versionCode"].(float64)),
				VersionName: packageMap["versionName"].(string),
				MinimumSDKVersion: packageMap["minimumSDKVersion"].(string),
				Permissions: permissions,
			},
		},
	}, nil
}
