# Copyright 2020 Huawei Technologies Co., Ltd.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#- name: Server is running on http port"
sed -i 's/value: \"0\"/value: \"1"/g' mepagent.yaml
sed -i 's/#volumeMounts/ volumeMounts/g' mepagent.yaml
sed -i 's/#- name/ - name/g' mepagent.yaml
sed -i 's/#mountPath/ mountPath/g' mepagent.yaml
sed -i 's/#readOnly/ readOnly/g' mepagent.yaml
sed -i 's/#volumes/ volumes/g' mepagent.yaml
sed -i 's/#- name/ - name/g' mepagent.yaml
sed -i 's/#secret/ secret/g' mepagent.yaml
sed -i 's/#secretName/ secretName/g' mepagent.yaml
