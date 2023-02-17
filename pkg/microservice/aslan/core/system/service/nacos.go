/*
Copyright 2023 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/koderover/zadig/pkg/microservice/aslan/core/common/repository/mongodb"
	"github.com/koderover/zadig/pkg/tool/nacos"
	"github.com/koderover/zadig/pkg/types"
)

func ListNacosNamespace(nacosID string, log *zap.SugaredLogger) ([]*types.NacosNamespace, error) {
	client, err := getNacosClient(nacosID)
	if err != nil {
		err = errors.Wrap(err, "fail to get nacos client")
		log.Error(err)
		return []*types.NacosNamespace{}, err
	}
	resp, err := client.ListNamespaces()
	if err != nil {
		err = errors.Wrap(err, "fail to list nacos namespace")
		log.Error(err)
		return []*types.NacosNamespace{}, err
	}
	return resp, nil
}

func ListNacosConfig(nacosID, namespaceID string, log *zap.SugaredLogger) ([]*types.NacosConfig, error) {
	client, err := getNacosClient(nacosID)
	if err != nil {
		err = errors.Wrap(err, "fail to get nacos client")
		log.Error(err)
		return []*types.NacosConfig{}, err
	}
	resp, err := client.ListConfigs(namespaceID)
	if err != nil {
		err = errors.Wrap(err, "fail to list nacos config")
		log.Error(err)
		return []*types.NacosConfig{}, err
	}
	return resp, nil
}

func getNacosClient(nacosID string) (*nacos.Client, error) {
	info, err := mongodb.NewConfigurationManagementColl().GetNacosByID(context.Background(), nacosID)
	if err != nil {
		return nil, errors.Wrap(err, "get nacos info")
	}
	return nacos.NewNacosClient(info.ServerAddress, info.UserName, info.Password)
}
