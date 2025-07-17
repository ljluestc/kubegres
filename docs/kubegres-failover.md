# Kubegres Failover with STONITH Support

## Overview

Kubegres implements automated failover capabilities to ensure high availability of your PostgreSQL clusters. When a primary instance fails, Kubegres automatically promotes a replica to become the new primary, minimizing downtime.

## STONITH Mechanism

STONITH (Shoot The Other Node In The Head) is a failover protection mechanism that ensures the old primary is fully terminated before promoting a replica. This helps prevent split-brain scenarios where multiple primaries might exist simultaneously.

### Why STONITH?

In distributed systems like PostgreSQL running on Kubernetes, network partitions or temporary connectivity issues might cause a scenario where the old primary is still running but can't communicate with the controller. If a new primary is promoted while the old one is still active, this can lead to data inconsistency and corruption.

STONITH ensures that:
1. The old primary is fully terminated
2. No conflicting writes can occur
3. Data integrity is maintained during failover

### Using STONITH in Kubegres

To enable STONITH for your Kubegres cluster, set the `enableSTONITH` flag in your Kubegres specification:
