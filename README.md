[![Build Status](https://travis-ci.org/RackHD/inservice.svg?branch=master)](https://travis-ci.org/RackHD/inservice)
[![Coverage Status](https://coveralls.io/repos/github/RackHD/inservice/badge.svg?branch=master)](https://coveralls.io/github/RackHD/inservice?branch=master)

# InService

## Agent

Compute, network, and storage devices all have different functionality which require different control surfaces and monitoring capabilities.  Many of the required features are not accessible via external interfaces and require services hosted on the host to make functionality accessible. Given the diverse sets of functionality across network, compute and storage systems, a plugin based architecture will allow for the reuse of common functionality amongst devices as well as provide a platform for new plugins to be included in the future.

The Agent is responsible for identifying what Plugins should be hosted on a device and will manage their lifecycle as child processes. It will provide functionality for the registration and advertisement of not only itself but the plugins it is hosting.  As such it will be responsible for implementing the advertisement portion of the SSDP protocol from https://tools.ietf.org/html/draft-cai-ssdp-v1-03. The advertisement will provide a URI for obtaining more information about the services (Plugins) available via the Agent.

Plugins should offer discrete functionality to the ecosystem which can be built upon by higher order services. They should generally be stateless themselves though they may manage the state of the device they are executing on. Any state provided should be recoverable by simply restarting the plugin or the hosting Agent. Dependencies on external state management should be avoided where possible.

###Overview

On startup the Agent will identify what Plugins should be hosted based on configuration and execute them as child processes. Each child process is assumed to be serving an RPC end point which implements start, status, and stop service methods. Once the Plugin process has been started the Agent will contact the Plugin on the RPC end point to notify it to start. In the event no errors occur the Agent will monitor the Plugin process and restart if it exits.

The Agent itself will host an RPC end point for managing Plugins which consists of a list of running plugins & start/status/stop service functionality for each service by it's name. When stopping a Plugin the Agent will call the Plugin RPC stop end point and allow the Plugin time to gracefully exit.  A status listing will also be made available.

In addition the Agent will advertise via SSDP and provide a URI for advertisement recipients to contact for more information about supported plugins.  The location URI will be hosted by the Agent and serve content from a local cache of data provided the plugins which register for advertisement with the Agent. As information changes the content provided by the location URI will change.

## Plugins
### Catalog Compute

RackHD (https://github.com/RackHD) is able to generate low level inventories of a resource's hardware components as part of its discovery process. Customers desire to have this cataloging process be made available while a server is running its primary operating system. InService Agent will execute this plugin periodically to generate a local equivalent of the hosts hardware catalog that can be consumed by infrastructure managers.

### LLDP

LLDP (IEEE standard 802.1AB-2009) is a network protocol for identifying devices physically connected to a hosts network port. It can be used to discover and describe network topology. The LLDP plugin is responsible for listening to LLDP advertisements on all of host's interfaces where the plugin is running and caching the data. The initial usage of this plugin will be targeted to InService running on a network switch to enable switch topology mappings.  The plugin may later be targeted to run on all compute nodes to provide mapping information from both sides of the network connection. The collected data will be made available via an API and higher order services can utilize the data to discovery devices and create a topology by combining LLDP data from multiple LLDP plugins.

[Agent]: https://github.com/RackHD/inservice/tree/master/agent
[Catalog Compute]: https://github.com/RackHD/inservice/tree/master/plugins/catalog-compute
[LLDP]: https://github.com/RackHD/inservice/tree/master/plugins/lldp

Contribute
----------

InService is a collection of libraries and applications housed at https://github.com/RackHD/inservice. The code for Inservice is written in Golang and makes use of Makefiles. It is available under the Apache 2.0 license (or compatible sublicences for library dependencies).

Code and bug submissions are handled on GitHub using the Issues tab for this repository above.

Community
---------

We also have a #InfraEnablers Slack channel: You can get an invite by requesting one at http://community.emccode.com.

Documentation
-------------

TODO:

Licensing
---------

Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

RackHD is a Trademark of EMC Corporation

Support
-------

Please file bugs and issues at the GitHub issues page. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
