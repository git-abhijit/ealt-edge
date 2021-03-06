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

package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/pkg/util"
	"github.com/apache/servicecomb-service-center/server/core/proto"

	meputil "mepserver/mp1/util"
)

// This type represents the general information of a MEC service info for registry/discovery
type ServiceInfo struct {
	SerInstanceId     string        `json:"serInstanceId,omitempty"`
	SerName           string        `json:"serName"`
	SerCategory       CategoryRef   `json:"serCategory,omitempty"`
	Version           string        `json:"version"`
	State             string        `json:"state"`
	TransportID       string        `json:"transportId,omitempty"`
	TransportInfo     TransportInfo `json:"transportInfo"`
	Serializer        string        `json:"serializer"`
	ScopeOfLocality   string        `json:"scopeOfLocality,omitempty"`
	ConsumedLocalOnly bool          `json:"consumedLocalOnly,omitempty"`
	IsLocal           bool          `json:"isLocal,omitempty"`
}

func (s *ServiceInfo) ToServiceRequest(req *proto.CreateServiceRequest) {
	if req == nil {
		log.Warn("create service request nil")
	}
	if req.Service == nil {
		req.Service = &proto.MicroService{}
	}
	req.Service.AppId = ""
	req.Service.ServiceName = s.SerName
	req.Service.Version = s.Version
	req.Service.Status = "UP"
	if s.State == "INACTIVE" {
		req.Service.Status = "DOWN"
	}
}

func (s *ServiceInfo) ToRegisterInstance(req *proto.RegisterInstanceRequest) {
	if req == nil {
		log.Warn("register instance request nil")
	}
	if req.Instance == nil {
		req.Instance = &proto.MicroServiceInstance{}
	}
	req.Instance.Properties = make(map[string]string, 5)
	req.Instance.Properties["serName"] = s.SerName
	s.serCategoryToProperties(req.Instance.Properties)
	req.Instance.Version = s.Version
	req.Instance.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	req.Instance.ModTimestamp = req.Instance.Timestamp

	req.Instance.Status = "UP"
	if s.State == "INACTIVE" {
		req.Instance.Status = "DOWN"
	}
	properties := req.Instance.Properties
	meputil.InfoToProperties(properties, "transportId", s.TransportID)
	meputil.InfoToProperties(properties, "serializer", s.Serializer)
	meputil.InfoToProperties(properties, "ScopeOfLocality", s.ScopeOfLocality)
	meputil.InfoToProperties(properties, "ConsumedLocalOnly", strconv.FormatBool(s.ConsumedLocalOnly))
	meputil.InfoToProperties(properties, "IsLocal", strconv.FormatBool(s.IsLocal))
	req.Instance.HostName = "default"
	var epType string
	req.Instance.Endpoints, epType = s.toEndpoints()
	req.Instance.Properties["endPointType"] = epType

	healthCheck := &proto.HealthCheck{
		Mode: proto.CHECK_BY_HEARTBEAT,
		Port: 0,
		Interval: math.MaxInt32 - 1,
		Times: 0,
		Url: "",
	}
	req.Instance.HealthCheck = healthCheck
	s.transportInfoToProperties(req.Instance.Properties)
}

func (s *ServiceInfo) toEndpoints() ([]string, string) {
	if len(s.TransportInfo.Endpoint.Addresses) != 0 {
		return s.TransportInfo.Endpoint.Uris, "uris"
	}
	endPoints := make([]string, 0, 1)
	if len(s.TransportInfo.Endpoint.Addresses) != 0 {

		for _, v := range s.TransportInfo.Endpoint.Addresses {
			addrDes := fmt.Sprintf("%s:%d", v.Host, v.Port)
			endPoints = append(endPoints, addrDes)
		}
		return endPoints, "addresses"
	}

	if s.TransportInfo.Endpoint.Alternative != nil {
		jsonBytes, err := json.Marshal(s.TransportInfo.Endpoint.Alternative)
		if err != nil {
			return nil, ""
		}
		jsonText := util.BytesToStringWithNoCopy(jsonBytes)
		endPoints = append(endPoints, jsonText)
		return endPoints, "alternative"
	}
	return nil, ""
}

func (s *ServiceInfo) transportInfoToProperties(properties map[string]string) {
	if properties == nil {
		return
	}
	meputil.InfoToProperties(properties, "transportInfo/id", s.TransportInfo.ID)
	meputil.InfoToProperties(properties, "transportInfo/name", s.TransportInfo.Name)
	meputil.InfoToProperties(properties, "transportInfo/description", s.TransportInfo.Description)
	meputil.InfoToProperties(properties, "transportInfo/type", string(s.TransportInfo.TransType))
	meputil.InfoToProperties(properties, "transportInfo/protocol", s.TransportInfo.Protocol)
	meputil.InfoToProperties(properties, "transportInfo/version", s.TransportInfo.Version)
	grantTypes := strings.Join(s.TransportInfo.Security.OAuth2Info.GrantTypes, "，")
	meputil.InfoToProperties(properties, "transportInfo/security/oAuth2Info/grantTypes", grantTypes)
	meputil.InfoToProperties(properties, "transportInfo/security/oAuth2Info/tokenEndpoint", s.TransportInfo.Security.OAuth2Info.TokenEndpoint)

}

func (s *ServiceInfo) serCategoryToProperties(properties map[string]string) {
	if properties == nil {
		return
	}
	meputil.InfoToProperties(properties, "serCategory/href", s.SerCategory.Href)
	meputil.InfoToProperties(properties, "serCategory/id", s.SerCategory.ID)
	meputil.InfoToProperties(properties, "serCategory/name", s.SerCategory.Name)
	meputil.InfoToProperties(properties, "serCategory/version", s.SerCategory.Version)
}

func (s *ServiceInfo) FromServiceInstance(inst *proto.MicroServiceInstance) {
	s.SerInstanceId = inst.ServiceId + inst.InstanceId
	s.serCategoryFromProperties(inst.Properties)
	s.Version = inst.Version
	s.State = "ACTIVE"
	if inst.Status == "DOWN" {
		s.State = "INACTIVE"
	}

	epType := "uris"
	if inst.Properties != nil {
		s.SerName = inst.Properties["serName"]
		s.TransportID = inst.Properties["transportId"]
		s.Serializer = inst.Properties["serializer"]
		epType = inst.Properties["endPointType"]
		s.ScopeOfLocality = inst.Properties["ScopeOfLocality"]
		var err error
		s.ConsumedLocalOnly, err = strconv.ParseBool(inst.Properties["ConsumedLocalOnly"])
		if err != nil {
			log.Warn("parse bool ConsumedLocalOnly fail")
		}
		s.IsLocal, err = strconv.ParseBool(inst.Properties["IsLocal"])
		if err != nil {
			log.Warn("parse bool ConsumedLocalOnly fail")
		}
	}
	s.fromEndpoints(inst.Endpoints, epType)
	s.transportInfoFromProperties(inst.Properties)
}

func (s *ServiceInfo) serCategoryFromProperties(properties map[string]string) {
	if properties == nil {
		return
	}
	s.SerCategory.Href = properties["serCategory/href"]
	s.SerCategory.ID = properties["serCategory/id"]
	s.SerCategory.Name = properties["serCategory/name"]
	s.SerCategory.Version = properties["serCategory/version"]
}

//Parse service endpoint info
func (s *ServiceInfo) fromEndpoints(uris []string, epType string) {
	if epType == "uris" {
		s.TransportInfo.Endpoint.Uris = uris
		return
	}
	if epType == "address" {

		s.TransportInfo.Endpoint.Addresses = make([]EndPointInfoAddress, 0, 1)
		for _, v := range uris {
			host, port := meputil.GetHostPort(v)
			tmp := EndPointInfoAddress{
				Host: host,
				Port: uint32(port),
			}
			s.TransportInfo.Endpoint.Addresses = append(s.TransportInfo.Endpoint.Addresses, tmp)
		}
	}
	if epType == "alternative" {
		jsonObj, err := meputil.JsonTextToObj(uris[0])
		if err != nil {
			s.TransportInfo.Endpoint.Alternative = jsonObj
		}
		return
	}
}

//Parse service transport info
func (s *ServiceInfo) transportInfoFromProperties(properties map[string]string) {
	if properties == nil {
		return
	}
	s.TransportInfo.ID = properties["transportInfo/id"]
	s.TransportInfo.Name = properties["transportInfo/name"]
	s.TransportInfo.Description = properties["transportInfo/description"]
	s.TransportInfo.TransType = TransportTypes(properties["transportInfo/type"])
	s.TransportInfo.Protocol = properties["transportInfo/protocol"]
	s.TransportInfo.Version = properties["transportInfo/version"]
	grantTypes := properties["transportInfo/security/oAuth2Info/grantTypes"]
	s.TransportInfo.Security.OAuth2Info.GrantTypes = strings.Split(grantTypes, ",")
	s.TransportInfo.Security.OAuth2Info.TokenEndpoint = properties["transportInfo/security/oAuth2Info/tokenEndpoint"]
}
