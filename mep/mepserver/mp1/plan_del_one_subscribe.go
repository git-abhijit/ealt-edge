/*
 * Copyright 2020 Huawei Technologies Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mp1

import (
	"context"
	"net/http"

	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/server/core/backend"
	"github.com/apache/servicecomb-service-center/server/core/proto"
	"github.com/apache/servicecomb-service-center/server/plugin/pkg/registry"

	"mepserver/mp1/arch/workspace"
)

type DelOneSubscribe struct {
	workspace.TaskBase
	R             *http.Request       `json:"r,in"`
	HttpErrInf    *proto.Response     `json:"httpErrInf,out"`
	Ctx           context.Context     `json:"ctx,in"`
	W             http.ResponseWriter `json:"w,in"`
	AppInstanceId string              `json:"appInstanceId,in"`
	SubscribeId   string              `json:"subscribeId,in"`
	HttpRsp       interface{}         `json:"httpRsp,out"`
}

func (t *DelOneSubscribe) OnRequest(data string) workspace.TaskCode {

	appInstanceId := t.AppInstanceId
	subscribeId := t.SubscribeId

	opts := []registry.PluginOp{
		registry.OpGet(registry.WithStrKey("/cse-sr/etsi/subscribe/" + appInstanceId + "/" + subscribeId)),
	}
	resp, errGet := backend.Registry().TxnWithCmp(context.Background(), opts, nil, nil)
	if errGet != nil {
		log.Errorf(errGet, "get subscription from etcd failed")
		t.SetFirstErrorCode(OperateDataWithEtcdErr, "get subscription from etch failed")
		return workspace.TaskFinish
	}

	if len(resp.Kvs) == 0 {
		log.Errorf(errGet, "subscription not exist")
		t.SetFirstErrorCode(SubscriptionNotFound, "subscription not exist")
		return workspace.TaskFinish
	}

	opts = []registry.PluginOp{
		registry.OpDel(registry.WithStrKey("/cse-sr/etsi/subscribe/" + appInstanceId + "/" + subscribeId)),
	}
	_, err := backend.Registry().TxnWithCmp(context.Background(), opts, nil, nil)
	if err != nil {
		log.Errorf(err, "delete subscription from etcd failed")
		t.SetFirstErrorCode(OperateDataWithEtcdErr, "delete subscription from etch failed")
		return workspace.TaskFinish
	}

	t.HttpRsp = ""
	return workspace.TaskFinish
}
