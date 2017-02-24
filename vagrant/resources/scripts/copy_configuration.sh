#!/bin/bash

set -e

. $HOME_DIR/.bash_profile

CONFIGS_LOCATION="/vagrant/resources/configs"

echo "Copying configs..."
mkdir -p /opt/swan/resources
mkdir -p /cache/ccache

cp $CONFIGS_LOCATION/docker.repo /etc/yum.repos.d/docker.repo
cp $CONFIGS_LOCATION/ccache.conf /etc/ccache.conf
cp $CONFIGS_LOCATION/fastestmirror.conf /etc/yum/pluginconf.d/fastestmirror.conf

cp $CONFIGS_LOCATION/cassandra.service /etc/systemd/system
cp $CONFIGS_LOCATION/keyspace.cql /opt/swan/resources
cp $CONFIGS_LOCATION/table.cql /opt/swan/resources

echo "Generate configs..."
sed -e 's,GOPATH,'"$GOPATH"',' $CONFIGS_LOCATION/snapteld.service_template > /etc/systemd/system/snapteld.service