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

package types

type NacosNamespace struct {
	NamespaceID    string `json:"namespace_id"`
	NamespacedName string `json:"namespace_name"`
}

type NacosConfig struct {
	DataID  string `bson:"data_id"                 json:"data_id"                 yaml:"data_id"`
	Group   string `bson:"group"                   json:"group"                   yaml:"group"`
	Desc    string `bson:"description,omitempty"   json:"description,omitempty"   yaml:"description,omitempty"`
	Format  string `bson:"format,omitempty"        json:"format,omitempty"        yaml:"format,omitempty"`
	Content string `bson:"content,omitempty"       json:"content,omitempty"       yaml:"content,omitempty"`
}
