/*
Copyright 2023 The KubeVela Authors.

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

package sharding

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildCache add shard-id label selector for given typed object
func BuildCache(scheme *runtime.Scheme, shardingObjects ...client.Object) cache.NewCacheFunc {
	return BuildCacheWithOptions(cache.Options{Scheme: scheme}, shardingObjects...)
}

// BuildCacheWithOptions add shard-id label selector to sharding objects with options
func BuildCacheWithOptions(opts cache.Options, shardingObjects ...client.Object) cache.NewCacheFunc {

	if EnableSharding {
		if opts.SelectorsByObject == nil {
			opts.SelectorsByObject = map[client.Object]cache.ObjectSelector{}
		}
		for _, obj := range shardingObjects {
			//master instance does not participate in sharding.
			if IsMaster() {
				opts.SelectorsByObject[obj] = cache.ObjectSelector{Label: labels.Everything()}
			} else {
				ls := labels.SelectorFromSet(map[string]string{LabelKubeVelaScheduledShardID: ShardID})
				opts.SelectorsByObject[obj] = cache.ObjectSelector{Label: ls}
			}
		}
	}
	return cache.BuilderWithOptions(opts)
}
