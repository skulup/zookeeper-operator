#
# Copyright 2020 Skulup Ltd, Open Collaborators
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM zookeeper:3.6.2

RUN apt-get update && apt-get install -y --no-install-recommends \
    lsof procps dnsutils \
    python \
    python-pip \
    python-setuptools \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN pip install zk-shell

RUN mkdir -p /zk && cp -r /apache-zookeeper-3.6.2-bin/* /zk
COPY deployments/docker/zookeeper/scripts /scripts
RUN chmod +x /scripts/* /zk/bin/*

ENTRYPOINT [ "/scripts/zkStart.sh" ]