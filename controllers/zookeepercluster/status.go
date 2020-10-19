/*
 * Copyright 2020 Skulup Ltd, Open Collaborators
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zookeepercluster

import (
	"context"
	"github.com/skulup/operator-helper/reconciler"
	"github.com/skulup/zookeeper-operator/api/v1alpha1"
	"github.com/skulup/zookeeper-operator/internal/zk_util"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ReconcileClusterStatus(ctx reconciler.Context, cluster *v1alpha1.ZookeeperCluster) (err error) {
	err = setZkMetaSizeCreated(ctx, cluster)
	return err
}

func setZkMetaSizeCreated(ctx reconciler.Context, cluster *v1alpha1.ZookeeperCluster) error {
	if !cluster.Status.ZkMetadata.SizeCreated {
		sts := &v1.StatefulSet{}
		return ctx.GetResource(types.NamespacedName{
			Name:      cluster.StatefulSetName(),
			Namespace: cluster.Namespace,
		}, sts,
			func() (err error) {
				if err = zk_util.UpdateZkClusterMetaSize(cluster); err == nil {
					cluster.Status.ZkMetadata.SizeCreated = true
					err = ctx.Client().Status().Update(context.TODO(), cluster)
				}
				return
			}, nil)
	}
	return nil
}
