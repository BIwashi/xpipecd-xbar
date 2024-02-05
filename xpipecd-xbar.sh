#!/usr/bin/env bash

# Metadata
# 
# <xbar.title>Pipecd xbar</xbar.title>
# <xbar.author>BIwashi</xbar.author>
# <xbar.author.github>BIwashi</xbar.author.github>
# <xbar.dependencies>go</xbar.dependencies>

# Variables 
#
# <xbar.var>string(PIPECD_API_KEY=""): Pipecd API Key.</xbar.var>
# <xbar.var>string(PIPECD_HOST="pipecd.jp:443"): Pipecd Host.</xbar.var>

./.artifacts/xpipecd-xbar pipectl --api-key=$PIPECD_API_KEY --host=$PIPECD_HOST
